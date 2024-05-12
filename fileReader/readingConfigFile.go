package fileReader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"
)

type RateLimit struct {
    Endpoint  string `json:"endpoint"`
    Burst     int    `json:"burst"`
    Sustained int    `json:"sustained"`
}

type ResponseBuffer struct {
    http.ResponseWriter
    Buffer bytes.Buffer
}

type RateLimitConfig struct {
    RateLimitsPerEndpoint []RateLimit `json:"rateLimitsPerEndpoint"`
}

type RateLimitInfo struct {
    LastRefill time.Time
    Sustained  int
    Remaining  int
    Burst      int
}
var limits = make(map[string]*RateLimitInfo)

func ReadingConfigFile() {
	data, err := os.ReadFile("config.json");
	now := time.Now()
	
	if err != nil {
		fmt.Println("Error reading file:", err)
		return;
	}
	var config RateLimitConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return;
	}

	fmt.Println("Rate Limits Configuration:")
	
	for _, limit := range config.RateLimitsPerEndpoint {
		fmt.Println(limit);
		limits[limit.Endpoint] =&RateLimitInfo{LastRefill:now,Sustained: limit.Sustained,Remaining: int(limit.Burst),Burst: int(limit.Burst) }
	}
}

func GetLimits()(map[string]*RateLimitInfo){
	return limits;
}
func RefillToken(limitInfo *RateLimitInfo) {
    now := time.Now()
    elapsed := now.Sub(limitInfo.LastRefill).Seconds()
    tokensToAdd := int(elapsed) * limitInfo.Sustained / 60
    if tokensToAdd > 0 {
        limitInfo.Remaining = int(math.Min(float64(limitInfo.Burst), float64(limitInfo.Remaining+tokensToAdd)))
        limitInfo.LastRefill = now
    }
}