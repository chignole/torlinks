package utils

import (
	"log"
	"math"
	"net/http"
	"time"
)

func TruncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

func PingHealthCheck(url string) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Println("[ERROR] Error sending ping:", err)
	}
	resp.Body.Close()
}

func CalculatePercentage(a float64, b float64) float64 {
	p := math.Floor((a/b*100)*100) / 100
	return p
}
