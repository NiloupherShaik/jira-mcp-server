package weather

type CurrentWeatherToolArgs struct {
	Location string `json:"location"`
}

type CurrentUVIndexToolargs struct {
	UVIndex int `json:"uv-index"`
}
