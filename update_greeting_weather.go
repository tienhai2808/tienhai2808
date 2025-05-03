package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60) 
	}

	now := time.Now().In(loc)
	hour := now.Hour()
	greeting := getGreeting(hour)

	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	city := "Da Nang"

	var weatherText string
	if apiKey != "" {
		data, err := getWeatherData(apiKey, city)
		if err != nil {
			log.Printf("Lỗi khi lấy dữ liệu thời tiết: %v", err)
			weatherText = "# 🌡️ Đà Nẵng: Không thể lấy dữ liệu\n"
		} else {
			weatherText = data
		}
	} else {
		weatherText = "# 🌡️ Đà Nẵng: Không thể lấy dữ liệu (API key không được cung cấp)\n"
	}

	newContent := []string{
		weatherText,
		fmt.Sprintf("### %s\n\n", greeting),
	}

	data, err := os.ReadFile("README.md")
	if err != nil {
		log.Fatalf("Không thể đọc file README.md: %v", err)
	}
	
	content := strings.Split(string(data), "\n")
	filteredContent := []string{}
	
	skip := false
	for _, line := range content {
		if shouldSkipLine(line) {
			skip = true
			continue
		}
		
		if skip && line == "" {
			skip = false
			continue
		}
		
		if !skip {
			filteredContent = append(filteredContent, line)
		}
	}
	
	finalContent := append(newContent, filteredContent...)

	err = os.WriteFile("README.md", []byte(strings.Join(finalContent, "\n")), 0644)
	if err != nil {
		log.Fatalf("Không thể ghi file README.md: %v", err)
	}
	
	fmt.Println("README.md đã được cập nhật!")
}