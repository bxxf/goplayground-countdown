package handlers

import (
	"bytes"
	"fmt"
	"image/gif"
	"log"
	"net/http"
	"time"

	"github.com/bxxf/goplayground-countdowns/cache"
	gm "github.com/bxxf/goplayground-countdowns/gif"
)

var gen *gm.Generator
var gca *cache.Cache

func Init(generator *gm.Generator, gifCache *cache.Cache) {
	gen = generator
	gca = gifCache
}

func HandleInvalidDate(w http.ResponseWriter) {
	http.Error(w, "Invalid date", http.StatusBadRequest)
}

func HandleCacheHit(w http.ResponseWriter, cacheKey string, gifG gif.GIF) {
	err := writeResponse(w, &gifG)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
	log.Printf("Writing Cached GIF to Response")
}

func HandleCacheMiss(w http.ResponseWriter, dateStr string, targetDate time.Time) {
	remainingDuration := time.Until(targetDate)
	placeholderGIF := generatePlaceholderGIF(remainingDuration)
	err := writeResponse(w, placeholderGIF)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
	log.Printf("Writing Placeholder GIF to Response")
	go generateAndCacheGIFs(dateStr, remainingDuration)
}

func writeResponse(w http.ResponseWriter, gifG *gif.GIF) error {
	buf := &bytes.Buffer{}
	err := gif.EncodeAll(buf, gifG)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}
	w.Header().Set("Content-Type", "image/gif")
	_, err = buf.WriteTo(w)
	return err
}

func generatePlaceholderGIF(remainingTime time.Duration) *gif.GIF {
	gimg, err := gen.GenerateGIF(remainingTime, true)
	if err != nil {
		log.Printf("Failed to generate placeholder GIF: %v", err)
		return nil
	}
	return gimg
}

func generateGIFAsync(remainingTime time.Duration) *gif.GIF {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in generateGIFAsync:", r)
		}
	}()
	gimg, err := gen.GenerateGIF(remainingTime, false)
	if err != nil {
		log.Printf("Failed to generate GIF: %v", err)
		return nil
	}
	return gimg
}

var nextHourCached = make(map[string]time.Time)
var nextHourGenerating = make(map[string]bool)

// pre generate GIFs for the next hour and cache them to improve response time
func generateAndCacheGIFs(dateStr string, remainingTime time.Duration) {
	if nextHourGenerating[dateStr] {
		log.Println("Cache is already generating", dateStr)
		return
	}
	nextHourGenerating[dateStr] = true
	startTime := time.Now()
	remainingTimeInHour := 60

	log.Printf("Generating GIFs for %d minutes", remainingTimeInHour)

	targetDate, _ := time.Parse("2006-01-02", dateStr)

	if nextHour, found := nextHourCached[dateStr]; found {
		if time.Now().Before(nextHour) {
			log.Println("Next hour already cached", dateStr, nextHour)
			return
		}
	}

	for i := remainingTimeInHour; i > 0; i-- {

		remaining := targetDate.Add(time.Minute*time.Duration(i)).Sub(time.Now()) - time.Hour
		remainingHours := int(remaining.Hours()) - 1
		remainingMinutes := int(remaining.Minutes()) % 60

		cacheKey := fmt.Sprintf("%s-%02dh%02dm", dateStr, remainingHours, remainingMinutes)

		_, found := gca.Fetch(cacheKey)

		if !found {
			gifG := generateGIFAsync(remaining)
			if gifG != nil {
				log.Println("Caching GIF for duration", cacheKey)
				gca.Store(cacheKey, *gifG)
			}
		} else {
			log.Println("GIF already cached for duration", cacheKey)
			continue
		}

		elapsedTime := time.Since(startTime)
		if elapsedTime > time.Second*30 {
			break
		}
	}
	nextHourCached[dateStr] = targetDate.Add(time.Hour)
	nextHourGenerating[dateStr] = false
}
