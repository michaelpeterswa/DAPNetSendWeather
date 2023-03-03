package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/mitchellh/mapstructure"
	geojson "github.com/paulmach/go.geojson"
	"go.uber.org/zap"
	"nw.codes/godapnet"
)

func getWeatherData(logger *zap.Logger, url string) []byte {
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

func parseWeatherData(logger *zap.Logger, data []byte) (*Properties, error) {
	feature, err := geojson.UnmarshalFeature(data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal feature: %w", err)
	}
	var properties *Properties
	err = mapstructure.Decode(feature.Properties, &properties)
	if err != nil {
		return nil, fmt.Errorf("could not decode properties: %w", err)
	}
	return properties, err
}

func sendCurrentForecast(logger *zap.Logger, f Forecast, sender *godapnet.Sender, settings DapnetSettings) error {
	msg := fmt.Sprintf("%s - %v%s - %s - Wind: %s %s", f.Name, f.Temperature, f.TemperatureUnit, f.ShortForecast, f.WindSpeed, f.WindDirection)
	callsigns := settings.CallsignNames
	txGps := settings.TransmitterGroupNames
	emerg := false

	// set configuration for message
	mc := godapnet.NewMessageConfig(sender.Callsign, godapnet.Alphapoc602RMaxMessageLength, callsigns, txGps, emerg)

	// send message
	err := sender.Send(msg, mc)
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}

	return nil
}
