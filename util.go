package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type WeatherResponse struct {
	Weather []struct {
		ID   int    `json:"id"`
		Main string `json:"main"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
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

func getWeatherData(apiKey string, city string) (string, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s&lang=vi", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API trả về status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return "", err
	}

	if len(weatherResp.Weather) == 0 {
		return "", fmt.Errorf("Không có dữ liệu thời tiết")
	}

	weatherID := weatherResp.Weather[0].ID
	currentTemp := weatherResp.Main.Temp
	weatherIcon := getWeatherIcon(weatherID)

	weatherText := fmt.Sprintf("# %s Đà Nẵng: %d°C\n", weatherIcon, int(currentTemp+0.5))
	return weatherText, nil
}

func shouldSkipLine(line string) bool {
	weatherIcons := []string{"# ⛈️", "# 🌦️", "# 🌧️", "# ❄️", "# 🌫️", "# ☀️", "# ☁️", "# 🌡️"}
	greetingIcons := []string{"# 🌅", "# 🍜", "# 🌞", "# 🌙", "# 🌃"}
	
	for _, icon := range append(weatherIcons, greetingIcons...) {
		if strings.HasPrefix(line, icon) {
			return true
		}
	}
	
	if strings.HasPrefix(line, "Thời tiết hiện tại ở") {
		return true
	}
	
	if strings.HasPrefix(line, "### ") {
		return true
	}
	
	return false
}
