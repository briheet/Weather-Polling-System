package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var pollInterval = time.Second * 5

const (
	endpoint = "https://api.open-meteo.com/v1/forecast" //?latitude=52.52&longitude=13.41&hourly=temperature_2m"
)

type WeatherData struct {
	Elevation float64        `json:"elevation"`
	Hourly    map[string]any `json:"hourly"`
}

type WPoller struct{}

func NewWPoller() *WPoller {
	return &WPoller{}
}

func (wp *WPoller) start() {
	fmt.Print("starting the WPoller")

	ticker := time.NewTicker(pollInterval)

	for {
		data, err := GetWeatherResults(52.52, 13.41)
		if err != nil {
			log.Fatal(err)
		}
		<-ticker.C
		data, err := GetWeatherResults(52.52, 13.41)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(data)
	}
}

func (wp *WPoller) HandleData(data *WeatherData) error {
	fmt.Println(data)
	return nil
}

func main() {
	wpoller := NewWPoller()
	wpoller.start()
}

func GetWeatherResults(lat, long float64) (*WeatherData, error) {
	uri := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m", endpoint, lat, long)

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
