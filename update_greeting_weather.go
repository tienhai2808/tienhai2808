package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// WeatherData represents the structure of the OpenWeatherMap API response
type WeatherData struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		ID int `json:"id"`
	} `json:"weather"`
}

func getWeatherIcon(weatherID int) string {
	switch {
	case weatherID >= 200 && weatherID <= 232:
		return "⛈️"
	case weatherID >= 300 && weatherID <= 321:
		return "🌦️"
	case weatherID >= 500 && weatherID <= 531:
		return "🌧️"
	case weatherID >= 600 && weatherID <= 622:
		return "❄️"
	case weatherID >= 701 && weatherID <= 781:
		return "🌫️"
	case weatherID == 800:
		return "☀️"
	case weatherID >= 801 && weatherID <= 804:
		return "☁️"
	default:
		return "🌡️"
	}
}

func main() {
	// Set timezone to Asia/Ho_Chi_Minh
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		fmt.Println("Error loading timezone:", err)
		return
	}
	now := time.Now().In(loc)
	hour := now.Hour()

	// Determine greeting based on time
	var greeting string
	switch {
	case hour >= 5 && hour < 11:
		greeting = "🌅 Chào buổi sáng! Hôm nay bạn đã code chưa?"
	case hour >= 11 && hour < 14:
		greeting = "🍜 Chào buổi trưa! Nghỉ ngơi một chút rồi code tiếp nào!"
	case hour >= 14 && hour < 18:
		greeting = "🌞 Chào buổi chiều! Hãy hoàn thành những dòng code cuối cùng!"
	case hour >= 18 && hour < 23:
		greeting = "🌙 Chào buổi tối! Push code xong thì đi ngủ sớm nhé!"
	default:
		greeting = "🌃 Khuya rồi! Nghỉ ngơi đi coder ơi!"
	}

	// Get weather data
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	city := "Da Nang"
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s&lang=vi", city, apiKey)

	var weatherText string
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		weatherText = "# 🌡️ Đà Nẵng: Không thể lấy dữ liệu\n"
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			weatherText = "# 🌡️ Đà Nẵng: Không thể lấy dữ liệu\n"
		} else {
			var data WeatherData
			if err := json.Unmarshal(body, &data); err != nil {
				weatherText = "# 🌡️ Đà Nẵng: Không thể lấy dữ liệu\n"
			} else {
				currentTemp := data.Main.Temp
				weatherID := data.Weather[0].ID
				weatherIcon := getWeatherIcon(weatherID)
				weatherText = fmt.Sprintf("# %s Đà Nẵng: %d°C\n", weatherIcon, int(currentTemp+0.5))
			}
		}
	}

	// Prepare new content
	newContent := []string{
		weatherText,
		fmt.Sprintf("### %s\n\n", greeting),
	}

	// Read existing README.md
	content, err := ioutil.ReadFile("README.md")
	if err != nil {
		fmt.Println("Error reading README.md:", err)
		return
	}
	lines := strings.Split(string(content), "\n")

	// Filter content
	var filteredContent []string
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
		if skip && strings.TrimSpace(line) == "" {
			skip = false
			continue
		}
		filteredContent = append(filteredContent, line)
	}

	// Combine new and filtered content
	finalContent := append(newContent, filteredContent...)

	// Write to README.md
	err = ioutil.WriteFile("README.md", []byte(strings.Join(finalContent, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing README.md:", err)
		return
	}

	fmt.Println("README.md đã được cập nhật!")
}