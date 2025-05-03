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

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		ID int `json:"id"`
	} `json:"weather"`
}

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

func getGreeting(hour int) string {
	switch {
	case hour >= 5 && hour < 11:
		return "🌅 Chào buổi sáng! Hôm nay bạn đã code chưa?"
	case hour >= 11 && hour < 14:
		return "🍜 Chào buổi trưa! Nghỉ ngơi một chút rồi code tiếp nào!"
	case hour >= 14 && hour < 18:
		return "🌞 Chào buổi chiều! Hãy hoàn thành những dòng code cuối cùng!"
	case hour >= 18 && hour < 23:
		return "🌙 Chào buổi tối! Push code xong thì đi ngủ sớm nhé!"
	default:
		return "🌃 Khuya rồi! Nghỉ ngơi đi coder ơi!"
	}
}

func main() {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.UTC
		fmt.Println("Error loading timezone, using UTC instead:", err)
	}

	now := time.Now().In(loc)
	hour := now.Hour()

	greeting := getGreeting(hour)

	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	city := "DaNang"
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s&lang=vi", city, apiKey)

	weatherText := ""

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			var weatherResp WeatherResponse
			err = json.Unmarshal(body, &weatherResp)
			
			if err == nil && len(weatherResp.Weather) > 0 {
				currentTemp := weatherResp.Main.Temp
				weatherID := weatherResp.Weather[0].ID
				weatherIcon := getWeatherIcon(weatherID)
				weatherText = fmt.Sprintf("# %s Đà Nẵng: %d°C\n", weatherIcon, int(currentTemp+0.5))
			}
		}
	}

	if weatherText == "" {
		weatherText = "# 🌡️ Đà Nẵng: Không thể lấy dữ liệu\n"
	}

	newContent := []string{
		weatherText,
		fmt.Sprintf("### %s\n\n", greeting),
	}

	readmeContent, err := ioutil.ReadFile("README.md")
	if err != nil {
		fmt.Println("Error reading README.md:", err)
		return
	}

	lines := strings.Split(string(readmeContent), "\n")
	filteredContent := []string{}

	skip := false
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if strings.HasPrefix(line, "# ⛈️") || strings.HasPrefix(line, "# 🌦️") || 
		   strings.HasPrefix(line, "# 🌧️") || strings.HasPrefix(line, "# ❄️") || 
		   strings.HasPrefix(line, "# 🌫️") || strings.HasPrefix(line, "# ☀️") || 
		   strings.HasPrefix(line, "# ☁️") || strings.HasPrefix(line, "# 🌡️") {
			skip = true

			continue
		}
		
		if skip && strings.HasPrefix(line, "### ") {
			if strings.Contains(line, "🌅") || strings.Contains(line, "🍜") || 
			   strings.Contains(line, "🌞") || strings.Contains(line, "🌙") || 
			   strings.Contains(line, "🌃") {
				continue
			}
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

	outputContent := strings.Join(finalContent, "\n")
	err = ioutil.WriteFile("README.md", []byte(outputContent), 0644)
	if err != nil {
		fmt.Println("Error writing to README.md:", err)
		return
	}

	fmt.Println("README.md đã được cập nhật!")
}