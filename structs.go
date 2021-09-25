package main

type DapnetSettings struct {
	CallsignNames         []string `yaml:"callsignNames"`
	TransmitterGroupNames []string `yaml:"transmitterGroupNames"`
}

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

type Properties struct {
	Updated           string     `mapstructure:"updated"`
	Units             string     `mapstructure:"units"`
	ForecastGenerator string     `mapstructure:"forecastGenerator"`
	GeneratedAt       string     `mapstructure:"generatedAt"`
	UpdateTime        string     `mapstructure:"updateTime"`
	ValidTimes        string     `mapstructure:"validTimes"`
	ElevationData     Elevation  `mapstructure:"elevation"`
	Periods           []Forecast `mapstructure:"periods"`
}

type Elevation struct {
	Value    int    `mapstructure:"value"`
	UnitCode string `mapstructure:"unitCode"`
}
