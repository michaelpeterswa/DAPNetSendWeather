package main

import (
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"nw.codes/godapnet"
)

func main() {
	var settings DapnetSettings

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	logger.Info("DAPNetSendWeather is initializing...")

	err = godotenv.Load("./config/.env")
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

	me := godapnet.NewSender(&http.Client{
		Timeout: 10 * time.Second,
	},
		godapnet.DAPNetURL,
		os.Getenv("CALLSIGN"),
		os.Getenv("DAPNET_USERNAME"),
		os.Getenv("DAPNET_PASSWORD"),
	)

	weatherCron := NewWeatherCron(settings, me, logger)

	c := cron.New()
	_, err = c.AddFunc(os.Getenv("CRON_STRING"), weatherCron.Run())
	if err != nil {
		logger.Fatal("failed to add cron function", zap.Error(err))
	}

	c.Start()
	select {}
}

type WeatherCron struct {
	settings DapnetSettings
	sender   *godapnet.Sender
	logger   *zap.Logger
}

func NewWeatherCron(settings DapnetSettings, sender *godapnet.Sender, logger *zap.Logger) *WeatherCron {
	return &WeatherCron{
		settings: settings,
		sender:   sender,
		logger:   logger,
	}
}

func (wc *WeatherCron) Run() func() {
	return func() {
		rawData := getWeatherData(wc.logger, os.Getenv("WEATHER_API_URL"))
		data, err := parseWeatherData(wc.logger, rawData)
		if err != nil {
			wc.logger.Error("Could not parse weather data", zap.Error(err))
			return
		}

		for _, forecast := range data.Periods {
			startTime, err := time.Parse(time.RFC3339, forecast.StartTime)
			if err != nil {
				wc.logger.Error("Could not parse startTime", zap.Error(err))
				continue
			}
			endTime, err := time.Parse(time.RFC3339, forecast.EndTime)
			if err != nil {
				wc.logger.Error("Could not parse endTime", zap.Error(err))
				continue
			}

			if startTime.Before(time.Now()) && endTime.After(time.Now()) {
				wc.logger.Info("Sending Forecast")
				err := sendCurrentForecast(wc.logger, forecast, wc.sender, wc.settings)
				if err != nil {
					wc.logger.Error("Could not send forecast", zap.Error(err))
					continue
				}
			}
		}
	}
}
