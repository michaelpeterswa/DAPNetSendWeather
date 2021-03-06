package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	dapnet "github.com/michaelpeterswa/godapnet"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func main() {
	var settings DapnetSettings

	err := godotenv.Load("./config/.env")
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	fileSettings, err := os.ReadFile("./config/settings.yaml")
	if err != nil {
		logger.Fatal("Error loading settings.yaml file")
	}

	err = yaml.Unmarshal(fileSettings, &settings)
	if err != nil {
		logger.Error("YAML failed to unmarshal to DapnetSettings", zap.Error(err))
	}

	me := dapnet.Sender{
		Callsign: os.Getenv("CALLSIGN"),
		Username: os.Getenv("DAPNET_USERNAME"),
		Password: os.Getenv("DAPNET_PASSWORD"),
	}

	c := cron.New()
	_, err = c.AddFunc("@every 2h0m0s", func() {
		rawData := getWeatherData(os.Getenv("WEATHER_API_URL"))
		data := parseWeatherData(rawData)

		for _, forecast := range data.Periods {
			startTime, err := time.Parse(time.RFC3339, forecast.StartTime)
			if err != nil {
				logger.Fatal("Could not parse startTime", zap.Error(err))
			}
			endTime, err := time.Parse(time.RFC3339, forecast.EndTime)
			if err != nil {
				logger.Fatal("Could not parse endTime", zap.Error(err))
			}

			if startTime.Before(time.Now()) && endTime.After(time.Now()) {
				logger.Info("Sending Forecast", zap.String("forecast", forecast.DetailedForecast))
				sendCurrentForecast(forecast, me, settings)
			}
		}
	})
	if err != nil {
		logger.Fatal("failed to add cron function", zap.Error(err))
	}

	c.Start()
	select {}
}
