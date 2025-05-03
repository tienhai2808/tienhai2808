package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// WeatherResponse represents the OpenWeatherMap API response structure
type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		ID int `json:"id"`
	} `json:"weather"`
}

// getWeatherIcon returns the appropriate emoji based on weather ID
func getWeatherIcon(weatherID int) string {
	if weatherID >= 200 && weatherID <= 232 {
		return "⛈️"
	} else if weatherID >= 300 && weatherID <= 321 {
		return "🌦️"
	} else if weatherID >= 500 && weatherID <= 531 {
		return "🌧️"
	} else if weatherID >= 600 && weatherID <= 622 {
		return "❄️"
	} else if weatherID >= 701 && weatherID <= 781 {
		return "🌫️"
	} else if weatherID == 800 {
		return "☀️"
	} else if weatherID >= 801 && weatherID <= 804 {
		return "☁️"
	} else {
		return "🌡️"
	}
}

func main() {
	// Print starting message
	log.Println("Starting weather README updater...")
	
	// Set timezone to Asia/Ho_Chi_Minh
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.UTC
		log.Println("Error loading timezone, using UTC instead:", err)
	}
	
	// Get current time in the specified timezone
	now := time.Now().In(loc)
	hour := now.Hour()

	// Determine greeting based on time of day
	var greeting string
	if hour >= 5 && hour < 11 {
		greeting = "🌅 Chào buổi sáng! Hôm nay bạn đã code chưa?"
	} else if hour >= 11 && hour < 14 {
		greeting = "🍜 Chào buổi trưa! Nghỉ ngơi một chút rồi code tiếp nào!"
	} else if hour >= 14 && hour < 18 {
		greeting = "🌞 Chào buổi chiều! Hãy hoàn thành những dòng code cuối cùng!"
	} else if hour >= 18 && hour < 23 {
		greeting = "🌙 Chào buổi tối! Push code xong thì đi ngủ sớm nhé!"
	} else {
		greeting = "🌃 Khuya rồi! Nghỉ ngơi đi coder ơi!"
	}

	// Get weather data
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	if apiKey == "" {
		log.Println("Warning: OPENWEATHERMAP_API_KEY environment variable is not set")
	}
	city := "Da Nang"
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s&lang=vi", city, apiKey)
	log.Println("Requesting weather data from:", url)

	weatherText := ""
	
	// Make HTTP request to OpenWeatherMap API
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
	} else {
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error status code: %d\n", resp.StatusCode)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
			} else {
				var weatherResp WeatherResponse
				err = json.Unmarshal(body, &weatherResp)
				
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
				} else if len(weatherResp.Weather) > 0 {
					currentTemp := weatherResp.Main.Temp
					weatherID := weatherResp.Weather[0].ID
					weatherIcon := getWeatherIcon(weatherID)
					weatherText = fmt.Sprintf("# %s Đà Nẵng: %d°C\n", weatherIcon, int(currentTemp+0.5))
					// Debug success
					fmt.Printf("Weather data fetched successfully: ID=%d, Temp=%.1f\n", weatherID, currentTemp)
				}
			}
		}
	}
	
	// If there was any error in the weather fetching process
	if weatherText == "" {
		weatherText = "# 🌡️ Đà Nẵng: Không thể lấy dữ liệu\n"
	}

	// Prepare new content
	newContent := []string{
		weatherText,
		fmt.Sprintf("### %s\n\n", greeting),
	}

	// Read existing README.md
	readmeContent, err := ioutil.ReadFile("README.md")
	if err != nil {
		fmt.Println("Error reading README.md:", err)
		return
	}

	// Split content into lines
	lines := strings.Split(string(readmeContent), "\n")
	filteredContent := []string{}
	
	// Filter out old weather and greeting information
	skip := false
	for _, line := range lines {
		if strings.HasPrefix(line, "# ⛈️") || strings.HasPrefix(line, "# 🌦️") || 
		   strings.HasPrefix(line, "# 🌧️") || strings.HasPrefix(line, "# ❄️") || 
		   strings.HasPrefix(line, "# 🌫️") || strings.HasPrefix(line, "# ☀️") || 
		   strings.HasPrefix(line, "# ☁️") || strings.HasPrefix(line, "# 🌡️") ||
		   strings.HasPrefix(line, "# 🌅") || strings.HasPrefix(line, "# 🍜") || 
		   strings.HasPrefix(line, "# 🌞") || strings.HasPrefix(line, "# 🌙") || 
		   strings.HasPrefix(line, "# 🌃") {
			skip = true
			continue
		}
		
		if skip && (strings.HasPrefix(line, "Thời tiết hiện tại ở") || 
			strings.HasPrefix(line, "# ⛈️") || strings.HasPrefix(line, "# 🌦️") || 
			strings.HasPrefix(line, "# 🌧️") || strings.HasPrefix(line, "# ❄️") || 
			strings.HasPrefix(line, "# 🌫️") || strings.HasPrefix(line, "# ☀️") || 
			strings.HasPrefix(line, "# ☁️") || strings.HasPrefix(line, "# 🌡️")) {
			continue
		}
		
		if skip && strings.HasPrefix(line, "### ") {
			continue
		}
		
		if skip && line == "" {
			skip = false
			continue
		}
		
		filteredContent = append(filteredContent, line)
	}

	// Combine new content with filtered content
	var finalContent []string
	finalContent = append(finalContent, newContent...)
	finalContent = append(finalContent, filteredContent...)

	// Join lines and write back to README.md
	outputContent := strings.Join(finalContent, "\n")
	err = ioutil.WriteFile("README.md", []byte(outputContent), 0644)
	if err != nil {
		fmt.Println("Error writing to README.md:", err)
		return
	}

	fmt.Println("README.md đã được cập nhật!")
}