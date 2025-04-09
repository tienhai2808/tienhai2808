import datetime
import requests
import os
import pytz

def get_weather_icon(weather_id):
  if 200 <= weather_id <= 232:
    return "⛈️"  
  elif 300 <= weather_id <= 321:
    return "🌦️"  
  elif 500 <= weather_id <= 531:
    return "🌧️"  
  elif 600 <= weather_id <= 622:
    return "❄️"  
  elif 701 <= weather_id <= 781:
    return "🌫️"  
  elif weather_id == 800:
    return "☀️" 
  elif 801 <= weather_id <= 804:
    return "☁️" 
  else:
    return "🌡️"  

tz = pytz.timezone('Asia/Ho_Chi_Minh')
now = datetime.datetime.now(tz)
hour = now.hour

if 5 <= hour < 11:
  greeting = "🌅 Chào buổi sáng! Hôm nay bạn đã code chưa?"
elif 11 <= hour < 14:
  greeting = "🍜 Chào buổi trưa! Nghỉ ngơi một chút rồi code tiếp nào!"
elif 14 <= hour < 18:
  greeting = "🌞 Chào buổi chiều! Hãy hoàn thành những dòng code cuối cùng!"
elif 18 <= hour < 23:
  greeting = "🌙 Chào buổi tối! Push code xong thì đi ngủ sớm nhé!"
else:
  greeting = "🌃 Khuya rồi! Nghỉ ngơi đi coder ơi!"

API_KEY = os.getenv("OPENWEATHERMAP_API_KEY")
CITY = "Da Nang"
URL = f"http://api.openweathermap.org/data/2.5/weather?q={CITY}&units=metric&appid={API_KEY}&lang=vi"

try:
  response = requests.get(URL)
  response.raise_for_status()  
  data = response.json()
  current_temp = data["main"]["temp"]
  weather_id = data["weather"][0]["id"]
  weather_icon = get_weather_icon(weather_id)
  weather_text = f"# {weather_icon} Đà Nẵng: {round(current_temp)}°C\n"
except Exception as e:
  weather_text = f"# 🌡️ Đà Nẵng: Không thể lấy dữ liệu\n"

new_content = [
  weather_text,
  f"### {greeting}\n\n"
]
with open("README.md", "r", encoding="utf-8") as f:
  content = f.readlines()

filtered_content = []
skip = False
for line in content:
  if (line.startswith("# ⛈️") or line.startswith("# 🌦️") or line.startswith("# 🌧️") or 
      line.startswith("# ❄️") or line.startswith("# 🌫️") or line.startswith("# ☀️") or 
      line.startswith("# ☁️") or line.startswith("# 🌡️") or 
      line.startswith("# 🌅") or line.startswith("# 🍜") or line.startswith("# 🌞") or 
      line.startswith("# 🌙") or line.startswith("# 🌃")):
    skip = True 
    continue
  if skip and (line.startswith("Thời tiết hiện tại ở") or 
               line.startswith("# ⛈️") or line.startswith("# 🌦️") or line.startswith("# 🌧️") or 
               line.startswith("# ❄️") or line.startswith("# 🌫️") or line.startswith("# ☀️") or 
               line.startswith("# ☁️") or line.startswith("# 🌡️")):
    continue
  if skip and line.startswith("### "):
    continue
  if skip and line.strip() == "":
    skip = False
    continue
  filtered_content.append(line)

final_content = new_content + filtered_content

with open("README.md", "w", encoding="utf-8") as f:
  f.writelines(final_content)

print("README.md đã được cập nhật!")