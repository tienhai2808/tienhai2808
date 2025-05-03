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
	// Thời gian theo Asia/Ho_Chi_Minh
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	now := time.Now().In(loc)
	hour := now.Hour()

	greeting := getGreeting(hour)

	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	city := "DaNang"
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s&lang=vi", city, apiKey)

	weatherText := ""
	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		var result map[string]interface{}
		json.Unmarshal(body, &result)

		if mainData, ok := result["main"].(map[string]interface{}); ok {
			temp := mainData["temp"].(float64)
			weatherList := result["weather"].([]interface{})
			weather := weatherList[0].(map[string]interface{})
			weatherID := int(weather["id"].(float64))

			icon := getWeatherIcon(weatherID)
			weatherText = fmt.Sprintf("# %s Đà Nẵng: %d°C\n", icon, int(temp+0.5))
		}
	}
	if weatherText == "" {
		weatherText = "# 🌡️ Đà Nẵng: Không thể lấy dữ liệu\n"
	}
	newContent := []string{
		weatherText,
		fmt.Sprintf("### %s\n\n", greeting),
	}

	// Đọc file README.md
	data, err := ioutil.ReadFile("README.md")
	if err != nil {
		fmt.Println("Không thể đọc README.md")
		return
	}
	lines := strings.Split(string(data), "\n")

	// Lọc nội dung cũ
	var filtered []string
	skip := false
	prefixes := []string{
		"# ⛈️", "# 🌦️", "# 🌧️", "# ❄️", "# 🌫️", "# ☀️", "# ☁️", "# 🌡️",
		"# 🌅", "# 🍜", "# 🌞", "# 🌙", "# 🌃",
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		startsWithPrefix := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(trimmed, prefix) {
				startsWithPrefix = true
				break
			}
		}
		if startsWithPrefix {
			skip = true
			continue
		}
		if skip {
			if strings.HasPrefix(trimmed, "### ") || trimmed == "" {
				skip = false
				continue
			} else {
				continue
			}
		}
		filtered = append(filtered, line)
	}

	finalContent := append(newContent, filtered...)
	output := strings.Join(finalContent, "\n")

	err = ioutil.WriteFile("README.md", []byte(output), 0644)
	if err != nil {
		fmt.Println("Không thể ghi README.md")
		return
	}

	fmt.Println("README.md đã được cập nhật!")
}
