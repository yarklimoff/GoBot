package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"encoding/json"
	"log"
)

type WeatherForecast struct {
	Fact struct {
		Temp int `json:"temp"`
		Humidity int `json:"humidity"`
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(telergam_api_key)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				btn := tgbotapi.NewKeyboardButtonLocation("Отправить геопозицию!")
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{btn})

				switch update.Message.Command() {
				case "start":
					msg.Text = "Привет, отправь свою геопозицию и получи прогноз погоды."
				default:
					msg.Text = "Попробуй /start"
				}
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			} else {
				if update.Message.Location != nil {
					lat := fmt.Sprintf("%f", update.Message.Location.Latitude)
					lon := fmt.Sprintf("%f", update.Message.Location.Longitude)

					url := "https://api.weather.yandex.ru/v2/informers?lat="+lat+"&lon="+lon
					req, _ := http.NewRequest("GET", url, nil)
					req.Header.Add("X-Yandex-API-Key", yandex_api_key)

					res, _ := http.DefaultClient.Do(req)
					defer res.Body.Close()
					body, _ := ioutil.ReadAll(res.Body)
					
					var weather WeatherForecast
					err := json.Unmarshal(body, &weather)
					if err == nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
						msg.Text = "Температура: " + fmt.Sprint(weather.Fact.Temp) + 
											 "\nВлажность: " + fmt.Sprint(weather.Fact.Humidity) + "%";
						if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
						}
					}
				}
			}
		} 
	}
}
