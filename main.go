package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	dapnet "github.com/michaelpeterswa/godapnet"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func main() {
	var settings DapnetSettings

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fileSettings, err := os.ReadFile("./settings.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(fileSettings, &settings)
	if err != nil {
		log.Printf("Failed to unmarshal: %s\n", err.Error())
	}

	me := dapnet.Sender{
		Callsign: os.Getenv("CALLSIGN"),
		Username: os.Getenv("DAPNET_USERNAME"),
		Password: os.Getenv("DAPNET_PASSWORD"),
	}

	rawData := getWeatherData(os.Getenv("WEATHER_API_URL"))
	data := parseWeatherData(rawData)

	for _, forecast := range data.Periods {
		startTime, err := time.Parse(time.RFC3339, forecast.StartTime)
		if err != nil {
			log.Panic(err.Error())
		}
		endTime, err := time.Parse(time.RFC3339, forecast.EndTime)
		if err != nil {
			log.Panic(err.Error())
		}

		if startTime.Before(time.Now()) && endTime.After(time.Now()) {
			logger.Info("Sending Forecast", zap.String("forecast", forecast.DetailedForecast))
			fmt.Println(me, settings)
			sendCurrentForecast(forecast, me, settings)
		}
	}
	logger.Info("Shutting Down...", zap.String("time", time.Now().Local().String()))
}
