package main

import (
	"fmt"
	"io"
	"net/http"

	dapnet "github.com/michaelpeterswa/godapnet"
	"github.com/mitchellh/mapstructure"
	geojson "github.com/paulmach/go.geojson"
	"go.uber.org/zap"
)

func getWeatherData(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		logger.Fatal("Could not get weather data from URL", zap.String("url", url), zap.Error(err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal("Could not parse response body", zap.Error(err))
	}
	return body
}

func parseWeatherData(data []byte) Properties {
	feature, err := geojson.UnmarshalFeature(data)
	if err != nil {
		logger.Fatal("Could not unmarshal GeoJSON feature", zap.Error(err))
	}
	var properties Properties
	err = mapstructure.Decode(feature.Properties, &properties)
	if err != nil {
		logger.Fatal("Mapstructure Decode Failed", zap.Error(err))
	}
	return properties
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
