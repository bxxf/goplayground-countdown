package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"

	"github.com/bxxf/goplayground-countdowns/cache"
	"github.com/bxxf/goplayground-countdowns/countdown"
	"github.com/bxxf/goplayground-countdowns/font"
	gp "github.com/bxxf/goplayground-countdowns/gif"
	"github.com/bxxf/goplayground-countdowns/handlers"
)

var gifCache *cache.Cache
var cacheMutex sync.RWMutex
var generator *gp.Generator

func init() {
	font.Init()
	gifCache = cache.NewCache(2000)
	generator = gp.NewGenerator(font.Face)
	handlers.Init(generator, gifCache)
}

func handler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	targetDate, err := time.Parse("2006-01-02", dateStr)

	if err != nil {
		handlers.HandleInvalidDate(w)
		return
	}

	remainingHours, remainingMinutes := countdown.GetRemainingTime(targetDate)
	cacheKey := fmt.Sprintf("%s-%02dh%02dm", dateStr, remainingHours, remainingMinutes)

	cacheMutex.RLock()
	gifG, found := gifCache.Fetch(cacheKey)
	cacheMutex.RUnlock()

	if found {
		log.Default().Println("Cache Hit", cacheKey)
		handlers.HandleCacheHit(w, cacheKey, gifG)
	} else {
		log.Default().Println("Cache Miss", cacheKey)
		handlers.HandleCacheMiss(w, dateStr, targetDate)
	}
}

func main() {
	http.HandleFunc("/countdown", handler)

	var portS string

	if port := os.Getenv("PORT"); port != "" {
		portS = ":" + port
	} else {
		portS = ":8080"
	}

	log.Println("Listening on port", portS)
	log.Fatal(http.ListenAndServe(portS, nil))

}
