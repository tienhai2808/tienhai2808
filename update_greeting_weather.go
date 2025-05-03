package main

import (
	"bufio"
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

	file, err := os.Open("README.md")
	if err != nil {
		fmt.Println("Error opening README.md:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var filteredContent []string
	skip := false
	for _, line := range lines {
		if strings.HasPrefix(line, "# ⛈️") || strings.HasPrefix(line, "# 🌦️") ||
			strings.HasPrefix(line, "# 🌧️") || strings.HasPrefix(line, "# ❄️") ||
			strings.HasPrefix(line, "# 🌫️") || strings.HasPrefix(line, "# ☀️") ||
			strings.HasPrefix(line, "# ☁️") || strings.HasPrefix(line, "# 🌡️") {
			skip = true
			continue
		}

		if skip && strings.HasPrefix(line, "### ") {
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

	var finalContent []string
	for _, line := range newContent {
		finalContent = append(finalContent, strings.TrimRight(line, "\n"))
	}
	finalContent = append(finalContent, filteredContent...)

	outputFile, err := os.Create("README.md")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	for i, line := range finalContent {
		fmt.Fprintln(writer, line)
		if i == len(finalContent)-1 && line == "" {
			continue
		}
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("README.md đã được cập nhật!")
}