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
	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	now := time.Now().In(location)
	hour := now.Hour()

	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	city := "Da Nang"
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s&lang=vi", city, apiKey)

	var weatherText string

	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		var data map[string]interface{}
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &data)

		main := data["main"].(map[string]interface{})
		weather := data["weather"].([]interface{})[0].(map[string]interface{})

		temp := int(main["temp"].(float64))
		weatherID := int(weather["id"].(float64))
		icon := getWeatherIcon(weatherID)

		weatherText = fmt.Sprintf("# %s Đà Nẵng: %d°C\n", icon, temp)
	} else {
		weatherText = "# 🌡️ Đà Nẵng: Không thể lấy dữ liệu\n"
	}

	greeting := fmt.Sprintf("### %s\n\n", getGreeting(hour))

	contentBytes, err := ioutil.ReadFile("README.md")
	if err != nil {
		fmt.Println("Không thể đọc README.md")
		return
	}
	lines := strings.Split(string(contentBytes), "\n")

	var filteredLines []string
	skip := false

	for _, line := range lines {
		if strings.HasPrefix(line, "# ⛈️") || strings.HasPrefix(line, "# 🌦️") || strings.HasPrefix(line, "# 🌧️") ||
			strings.HasPrefix(line, "# ❄️") || strings.HasPrefix(line, "# 🌫️") || strings.HasPrefix(line, "# ☀️") ||
			strings.HasPrefix(line, "# ☁️") || strings.HasPrefix(line, "# 🌡️") ||
			strings.HasPrefix(line, "### 🌅") || strings.HasPrefix(line, "### 🍜") || strings.HasPrefix(line, "### 🌞") ||
			strings.HasPrefix(line, "### 🌙") || strings.HasPrefix(line, "### 🌃") {
			skip = true
			continue
		}
		if skip && strings.TrimSpace(line) == "" {
			skip = false
			continue
		}
		if !skip {
			filteredLines = append(filteredLines, line)
		}
	}

	newLines := []string{weatherText, greeting}
	finalContent := append(newLines, filteredLines...)
	output := strings.Join(finalContent, "\n")

	err = ioutil.WriteFile("README.md", []byte(output), 0644)
	if err != nil {
		fmt.Println("Không thể ghi README.md")
	} else {
		fmt.Println("README.md đã được cập nhật!")
	}
}
