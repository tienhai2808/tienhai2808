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
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	now := time.Now().In(loc)
	hour := now.Hour()

	greeting := getGreeting(hour)

	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	city := "Da Nang"
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s&lang=vi", city, apiKey)

	weatherText := ""

	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		var weatherResp WeatherResponse
		if err := json.Unmarshal(body, &weatherResp); err == nil && len(weatherResp.Weather) > 0 {
			temp := int(weatherResp.Main.Temp + 0.5)
			icon := getWeatherIcon(weatherResp.Weather[0].ID)
			weatherText = fmt.Sprintf("# %s Đà Nẵng: %d°C\n", icon, temp)
		}
	}

	if weatherText == "" {
		weatherText = "# 🌡️ Đà Nẵng: Không thể lấy dữ liệu\n"
	}

	newContent := []string{
		weatherText,
		fmt.Sprintf("### %s\n", greeting),
		"",
	}

	readme, err := ioutil.ReadFile("README.md")
	if err != nil {
		fmt.Println("Không đọc được README.md:", err)
		return
	}

	lines := strings.Split(string(readme), "\n")
	filtered := []string{}
	skip := false

	for _, line := range lines {
		if strings.HasPrefix(line, "# ⛈️") || strings.HasPrefix(line, "# 🌦️") ||
			strings.HasPrefix(line, "# 🌧️") || strings.HasPrefix(line, "# ❄️") ||
			strings.HasPrefix(line, "# 🌫️") || strings.HasPrefix(line, "# ☀️") ||
			strings.HasPrefix(line, "# ☁️") || strings.HasPrefix(line, "# 🌡️") ||
			strings.HasPrefix(line, "# 🌅") || strings.HasPrefix(line, "# 🍜") ||
			strings.HasPrefix(line, "# 🌞") || strings.HasPrefix(line, "# 🌙") ||
			strings.HasPrefix(line, "# 🌃") || strings.HasPrefix(line, "### ") {
			skip = true
			continue
		}
		if skip && strings.TrimSpace(line) == "" {
			skip = false
			continue
		}
		if !skip {
			filtered = append(filtered, line)
		}
	}

	final := append(newContent, filtered...)
	output := strings.Join(final, "\n")

	err = ioutil.WriteFile("README.md", []byte(output), 0644)
	if err != nil {
		fmt.Println("Không thể ghi file:", err)
		return
	}

	fmt.Println("✅ README.md đã được cập nhật!")
}
