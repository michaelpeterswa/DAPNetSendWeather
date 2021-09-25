package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWeatherData(t *testing.T) {
	tests := []struct {
		name        string
		weatherData []byte
		result      Properties
	}{
		{
			weatherData: []byte(`{ "properties":{ "updated":"a", "units":"a", "forecastGenerator":"a", "generatedAt":"a", "updateTime":"a", "validTimes":"a", "elevation":{ "value":1, "unitCode":"a" }, "periods":[ { "number":1, "name":"Tonight", "startTime":"2021-09-24T20:00:00-07:00", "endTime":"2021-09-25T06:00:00-07:00", "isDaytime":false, "temperature":53, "temperatureUnit":"F", "temperatureTrend":null, "windSpeed":"6 mph", "windDirection":"N", "icon":"https://api.weather.gov/icons/land/night/few?size=medium", "shortForecast":"Mostly Clear", "detailedForecast":"Mostly clear, with a low around 53. North wind around 6 mph." } ] } }`),
			result: Properties{
				Updated:           "a",
				Units:             "a",
				ForecastGenerator: "a",
				GeneratedAt:       "a",
				UpdateTime:        "a",
				ValidTimes:        "a",
				ElevationData: Elevation{
					Value:    1,
					UnitCode: "a",
				},
				Periods: []Forecast{
					{
						Number:           1,
						Name:             "Tonight",
						StartTime:        "2021-09-24T20:00:00-07:00",
						EndTime:          "2021-09-25T06:00:00-07:00",
						IsDaytime:        false,
						Temperature:      53,
						TemperatureUnit:  "F",
						TemperatureTrend: nil,
						WindSpeed:        "6 mph",
						WindDirection:    "N",
						Icon:             "https://api.weather.gov/icons/land/night/few?size=medium",
						ShortForecast:    "Mostly Clear",
						DetailedForecast: "Mostly clear, with a low around 53. North wind around 6 mph.",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			props := parseWeatherData(tc.weatherData)
			assert.Equal(t, tc.result, props)
		})
	}
}
