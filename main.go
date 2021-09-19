package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/joho/godotenv"
	dapnet "github.com/michaelpeterswa/godapnet"
	"github.com/mitchellh/mapstructure"
	geojson "github.com/paulmach/go.geojson"
	"gopkg.in/yaml.v3"
)

type DapnetSettings struct {
	CallsignNames         []string `yaml:"callsignNames"`
	TransmitterGroupNames []string `yaml:"transmitterGroupNames"`
}

func getWeatherData(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return body
}

func parseWeatherData(data []byte) map[string]interface{} {
	feature, err := geojson.UnmarshalFeature(data)
	if err != nil {
		log.Println(err)
	}
	return feature.Properties
}

func sendCurrentForecast(f Forecast, sender dapnet.Sender, settings DapnetSettings) {
	msg := fmt.Sprintf("%s - %v%s - %s - Wind: %s %s", f.Name, f.Temperature, f.TemperatureUnit, f.ShortForecast, f.WindSpeed, f.WindDirection)
	callsigns := settings.CallsignNames
	txGps := settings.TransmitterGroupNames
	emerg := false

	messages := dapnet.CreateMessage(sender.Callsign, msg, callsigns, txGps, emerg)
	payloads := dapnet.GeneratePayload(messages)
	dapnet.SendMessage(payloads, sender.Username, sender.Password)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fileSettings, err := os.ReadFile("./settings.yaml")
	if err != nil {
		panic(err)
	}

	var settings DapnetSettings
	err = yaml.Unmarshal(fileSettings, &settings)
	if err != nil {
		log.Printf("Failed to unmarshal: %s\n", err.Error())
	}

	me := dapnet.Sender{
		Callsign: os.Getenv("CALLSIGN"),
		Username: os.Getenv("DAPNET_USERNAME"),
		Password: os.Getenv("DAPNET_PASSWORD"),
	}

	var nextForecastMap map[string]interface{}
	rawData := getWeatherData(os.Getenv("WEATHER_API_URL"))
	data := parseWeatherData(rawData)
	periods := data["periods"]
	var allForecasts []Forecast
	if reflect.TypeOf(periods).Kind() == reflect.Slice {
		forecasts := reflect.ValueOf(periods)
		for i := 0; i < forecasts.Len(); i++ {
			var fc Forecast
			nextForecast := forecasts.Index(i)
			nextForecastInterface := nextForecast.Interface()
			nextForecastMap = nextForecastInterface.(map[string]interface{})
			err := mapstructure.Decode(nextForecastMap, &fc)
			if err != nil {
				log.Panic(err.Error())
			}
			allForecasts = append(allForecasts, fc)
		}
	}

	for _, forecast := range allForecasts {
		startTime, err := time.Parse(time.RFC3339, forecast.StartTime)
		if err != nil {
			log.Panic(err.Error())
		}
		endTime, err := time.Parse(time.RFC3339, forecast.EndTime)
		if err != nil {
			log.Panic(err.Error())
		}

		if startTime.Before(time.Now()) && endTime.After(time.Now()) {
			sendCurrentForecast(forecast, me, settings)
		}
	}
}
