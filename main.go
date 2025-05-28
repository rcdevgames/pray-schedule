package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReverseGeocodeResponse struct {
	City        string `json:"city"`
	CountryName string `json:"countryName"`
}

type PrayerTimesResponse struct {
	Data struct {
		Timings map[string]string `json:"timings"`
	} `json:"data"`
}

func main() {
	r := gin.Default()

	r.GET("/prayer-times", func(c *gin.Context) {
		lat := c.Query("lat")
		lon := c.Query("lon")

		if lat == "" || lon == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lon are required"})
			return
		}

		// Step 1: Reverse Geocoding
		geoURL := fmt.Sprintf("https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=%s&longitude=%s&localityLanguage=en", lat, lon)
		resp, err := http.Get(geoURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to call reverse geocode API"})
			return
		}
		defer resp.Body.Close()

		var geoData ReverseGeocodeResponse
		if err := json.NewDecoder(resp.Body).Decode(&geoData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse reverse geocode response"})
			return
		}

		// Step 2: Get Prayer Times
		prayerURL := fmt.Sprintf("https://api.aladhan.com/v1/timingsByCity?city=%s&country=%s", geoData.City, geoData.CountryName)
		prayerResp, err := http.Get(prayerURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to call prayer times API"})
			return
		}
		defer prayerResp.Body.Close()

		var prayerData PrayerTimesResponse
		if err := json.NewDecoder(prayerResp.Body).Decode(&prayerData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse prayer times response"})
			return
		}

		// Final response
		c.JSON(http.StatusOK, gin.H{
			"city":    geoData.City,
			"country": geoData.CountryName,
			"times":   prayerData.Data.Timings,
		})
	})

	if err := r.Run(":5555"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
