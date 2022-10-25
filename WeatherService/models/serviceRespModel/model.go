package serviceRespModel

type CityWeather struct {
	Name       string  `json:"name"`
	Country    string  `json:"country"`
	AvTemp     float64 `json:"av_temp"`
	WeatherArr []Weather
}

type Weather struct {
	Temp  float64 `json:"temp"`
	Dt    int     `json:"dt"`
	DtTxt string  `json:"dt_txt"`
}
