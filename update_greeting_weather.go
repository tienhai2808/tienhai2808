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
	case 200 <= weatherID && weatherID <= 232:
		return "⛈️"
	case 300 <= weatherID && weatherID <= 321:
		return "🌦️"
	case 500 <= weatherID && weatherID <= 531:
		return "🌧️"
	case 600 <= weatherID && weatherID <= 622:
		return "❄️"
	case 701 <= weatherID && weatherID <= 781:
		return "🌫️"
	case weatherID == 800:
		return "☀️"
	case 801 <= weatherID && weatherID <= 804:
		return "☁️"
	default:
		return "🌡️"
	}
}

func getGreeting(hour int) string {
	switch {
	case 5 <= hour && hour < 11:
		return "🌅 Chào buổi sáng! Hôm nay bạn đã code chưa?"
	case 11 <= hour && hour < 14:
		return "🍜 Chào buổi trưa! Nghỉ ngơi một chút rồi code tiếp nào!"
	case 14 <= hour && hour < 18:
		return "🌞 Chào buổi chiều! Hãy hoàn thành những dòng code cuối cùng!"
	case 18 <= hour && hour < 23:
		return "🌙 Chào buổi tối! Push code xong thì đi ngủ sớm nhé!"
	default:
		return "🌃 Khuya rồi! Nghỉ ngơi đi coder ơi!"
	}
}

func main() {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	now := time.Now().In(loc)
	hour := now.Hour()
	greeting := getGreeting(hour)

	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	city := "DaNang"
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s&lang=vi", city, apiKey)

	var weatherText string
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		var data map[string]interface{}
		json.Unmarshal(body, &data)

		mainData := data["main"].(map[string]interface{})
		weatherArray := data["weather"].([]interface{})
		weather := weatherArray[0].(map[string]interface{})

		temp := int(mainData["temp"].(float64) + 0.5)
		weatherID := int(weather["id"].(float64))
		icon := getWeatherIcon(weatherID)
		weatherText = fmt.Sprintf("# %s Đà Nẵng: %d°C\n", icon, temp)
	} else {
		weatherText = "# 🌡️ Đà Nẵng: Không thể lấy dữ liệu\n"
	}

	newContent := []string{
		weatherText,
		fmt.Sprintf("### %s\n\n", greeting),
	}

	input, err := ioutil.ReadFile("README.md")
	if err != nil {
		fmt.Println("Không thể đọc file README.md")
		return
	}
	lines := strings.Split(string(input), "\n")

	filteredContent := []string{}
	skip := false

	for _, line := range lines {
		if strings.HasPrefix(line, "# ⛈️") || strings.HasPrefix(line, "# 🌦️") || strings.HasPrefix(line, "# 🌧️") ||
			strings.HasPrefix(line, "# ❄️") || strings.HasPrefix(line, "# 🌫️") || strings.HasPrefix(line, "# ☀️") ||
			strings.HasPrefix(line, "# ☁️") || strings.HasPrefix(line, "# 🌡️") ||
			strings.HasPrefix(line, "# 🌅") || strings.HasPrefix(line, "# 🍜") || strings.HasPrefix(line, "# 🌞") ||
			strings.HasPrefix(line, "# 🌙") || strings.HasPrefix(line, "# 🌃") {
			skip = true
			continue
		}
		if skip && (strings.HasPrefix(line, "Thời tiết hiện tại ở") ||
			strings.HasPrefix(line, "# ⛈️") || strings.HasPrefix(line, "# 🌦️") || strings.HasPrefix(line, "# 🌧️") ||
			strings.HasPrefix(line, "# ❄️") || strings.HasPrefix(line, "# 🌫️") || strings.HasPrefix(line, "# ☀️") ||
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

	finalContent := append(newContent, filteredContent...)
	output := strings.Join(finalContent, "\n")

	err = ioutil.WriteFile("README.md", []byte(output), 0644)
	if err != nil {
		fmt.Println("Không thể ghi file README.md")
		return
	}

	fmt.Println("README.md đã được cập nhật!")
}
