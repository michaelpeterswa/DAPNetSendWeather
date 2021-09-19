package main

type Forecast struct {
	Number           int     `mapstructure:"number"`
	Name             string  `mapstructure:"name"`
	StartTime        string  `mapstructure:"startTime"`
	EndTime          string  `mapstructure:"endTime"`
	IsDaytime        bool    `mapstructure:"isDaytime"`
	Temperature      int     `mapstructure:"temperature"`
	TemperatureUnit  string  `mapstructure:"temperatureUnit"`
	TemperatureTrend *string `mapstructure:"temperatureTrend"`
	WindSpeed        string  `mapstructure:"windSpeed"`
	WindDirection    string  `mapstructure:"windDirection"`
	Icon             string  `mapstructure:"icon"`
	ShortForecast    string  `mapstructure:"shortForecast"`
	DetailedForecast string  `mapstructure:"detailedForecast"`
}
