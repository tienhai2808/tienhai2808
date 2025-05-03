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
	city := "DaNang"
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
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		
		// Kiểm tra các dòng thời tiết/thời gian để bắt đầu bỏ qua
		if strings.HasPrefix(line, "# ⛈️") || strings.HasPrefix(line, "# 🌦️") || 
		   strings.HasPrefix(line, "# 🌧️") || strings.HasPrefix(line, "# ❄️") || 
		   strings.HasPrefix(line, "# 🌫️") || strings.HasPrefix(line, "# ☀️") || 
		   strings.HasPrefix(line, "# ☁️") || strings.HasPrefix(line, "# 🌡️") {
			skip = true
			
			// Bỏ qua dòng hiện tại và tiếp tục lặp
			continue
		}
		
		// Nếu đang trong chế độ bỏ qua và gặp dòng greeting
		if skip && strings.HasPrefix(line, "### ") {
			// Kiểm tra nếu là dòng greeting có emoji thời gian
			if strings.Contains(line, "🌅") || strings.Contains(line, "🍜") || 
			   strings.Contains(line, "🌞") || strings.Contains(line, "🌙") || 
			   strings.Contains(line, "🌃") {
				continue
			}
		}
		
		// Thoát chế độ bỏ qua khi gặp dòng trống sau các phần cần bỏ qua
		if skip && line == "" {
			skip = false
			continue  // Bỏ qua dòng trống này
		}
		
		// Nếu không trong chế độ bỏ qua, thêm dòng vào nội dung đã lọc
		if !skip {
			filteredContent = append(filteredContent, line)
		}
	}

	// Kết hợp nội dung mới và nội dung đã lọc
	finalContent := append(newContent, filteredContent...)

	outputContent := strings.Join(finalContent, "\n")
	err = ioutil.WriteFile("README.md", []byte(outputContent), 0644)
	if err != nil {
		fmt.Println("Error writing to README.md:", err)
		return
	}

	fmt.Println("README.md đã được cập nhật!")
}