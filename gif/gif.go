package gif

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"log"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
)

const (
	W          = 900
	H          = 250
	Minutes    = 60
	Loop       = Minutes
	ImageDelay = 6000
)

var palette = color.Palette{color.Black, color.White}

type Generator struct {
	font           *truetype.Font
	fontSize       float64
	dpi            float64
	minSecWidths   [60]float64
	separatorWidth float64
}

func NewGenerator(font *truetype.Font) *Generator {
	g := &Generator{
		font:     font,
		fontSize: W / 7,
		dpi:      72,
	}

	go g.calculateWidths()
	return g
}

func (g *Generator) calculateWidths() {
	opts := truetype.Options{
		Size: g.fontSize / 7,
		DPI:  g.dpi,
	}

	face := truetype.NewFace(g.font, &opts)

	// calculate widths for mins and secs
	for i := 0; i < 60; i++ {
		str := fmt.Sprintf("%02d", i)
		var width float64
		for _, r := range str {
			advance, ok := face.GlyphAdvance(r)
			if !ok {
				log.Println("Error while calculating string width:", str)
				continue
			}
			width += float64(advance >> 6) // convert from 26.6 fixed-point notation
		}
		g.minSecWidths[i] = width
	}

	// calculate width of separator
	var separatorWidth float64
	for _, r := range " : " {
		advance, _ := face.GlyphAdvance(r)
		separatorWidth += float64(advance >> 6)
	}
	g.separatorWidth = separatorWidth
}

func (g *Generator) drawCountdown(dc *gg.Context, remainingTime time.Duration, isPlaceholder bool) {
	times := g.getTimes(remainingTime)

	dc.SetRGB(0, 0, 0)
	dc.Fill()

	face := truetype.NewFace(g.font, &truetype.Options{Size: g.fontSize, DPI: g.dpi})
	dc.SetFontFace(face)
	sw, _ := dc.MeasureString(times.countdownStr)

	startX := (float64(dc.Width()) - sw) * 0.85
	cacheWidthAddon := g.getCacheWidthAddon(isPlaceholder)
	y := g.getYPosition(isPlaceholder)

	g.drawStringWithPositionUpdate(dc, &startX, times.dayStr, y, times.days, cacheWidthAddon)
	g.drawStringWithPositionUpdate(dc, &startX, " : ", y, 0, cacheWidthAddon)
	g.drawStringWithPositionUpdate(dc, &startX, times.hourStr, y, times.hours, cacheWidthAddon)
	g.drawStringWithPositionUpdate(dc, &startX, " : ", y, 0, cacheWidthAddon)
	g.drawStringWithPositionUpdate(dc, &startX, times.minuteStr, y, times.minutes, cacheWidthAddon)
}

func (g *Generator) getTimes(remainingTime time.Duration) *times {
	days := int(remainingTime.Hours() / 24)
	hours := int(remainingTime.Hours()) % 24
	minutes := int(remainingTime.Minutes()) % 60

	dayStr, hourStr, minuteStr := fmt.Sprintf("%02d", days), fmt.Sprintf("%02d", hours), fmt.Sprintf("%02d", minutes)

	if remainingTime.Milliseconds() < 0 {
		dayStr, hourStr, minuteStr = "00", "00", "00"
	}

	return &times{dayStr, hourStr, minuteStr, days, hours, minutes, dayStr + " : " + hourStr + " : " + minuteStr}
}

func (g *Generator) getCacheWidthAddon(isPlaceholder bool) float64 {
	if isPlaceholder {
		return 1.7
	}
	return 5.2
}

func (g *Generator) getYPosition(isPlaceholder bool) float64 {
	if isPlaceholder {
		return H/3/2 + g.fontSize/3
	}
	return H/2 + g.fontSize/3
}

func (g *Generator) drawStringWithPositionUpdate(dc *gg.Context, startX *float64, str string, y float64, value int, cacheWidthAddon float64) {
	g.drawString(dc, str, *startX, y)
	var sw float64
	if str == " : " { // If the string is a separator
		sw = g.separatorWidth * cacheWidthAddon
	} else {
		sw = g.getUpdatedWidth(value, cacheWidthAddon)
	}
	*startX += sw
}

func (g *Generator) getUpdatedWidth(value int, cacheWidthAddon float64) float64 {

	return g.minSecWidths[value] * cacheWidthAddon
}

func (g *Generator) drawString(dc *gg.Context, s string, x, y float64) {
	dc.DrawString(s, x, y)
}

func (g *Generator) GenerateGIF(remainingTime time.Duration, isPlaceholder bool) (*gif.GIF, error) {
	gifG := &gif.GIF{}

	for i := 0; i < Loop; i++ {
		remainingFrameTime := remainingTime - time.Duration(i)*time.Minute - time.Hour

		_, _, finalImg := g.createImage(remainingFrameTime, isPlaceholder)

		img := image.NewPaletted(image.Rect(0, 0, finalImg.Bounds().Dx(), finalImg.Bounds().Dy()), palette)
		draw.Draw(img, img.Bounds(), finalImg, image.Point{}, draw.Src)

		gifG.Image = append(gifG.Image, img)
		gifG.Delay = append(gifG.Delay, ImageDelay)

		if remainingFrameTime-time.Minute <= 0 {
			break
		}
	}

	return gifG, nil
}

func (g *Generator) createImage(remainingFrameTime time.Duration, isPlaceholder bool) (int, int, image.Image) {
	var width, height int
	if isPlaceholder {
		width, height = W/3, H/3
	} else {
		width, height = W, H
	}

	g.fontSize = float64(width) / 7

	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	g.drawCountdown(dc, remainingFrameTime, isPlaceholder)

	if !isPlaceholder {
		return width, height, resize.Resize(W/3, H/3, dc.Image(), resize.Lanczos3)
	}
	return width, height, dc.Image()
}

type times struct {
	dayStr, hourStr, minuteStr string
	days, hours, minutes       int
	countdownStr               string
}
