package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	poller    = time.Second * 2
	endpoller = time.Minute * 1
)

const endpoint = "https://api.open-meteo.com/v1/forecast" //?latitude=52.52&longitude=13.41&hourly=temperature_2m"

type Sender interface {
	Send(*WeatherData) error
}

type WeatherData struct {
	Elevation float64        `json:"elevation"`
	Hourly    map[string]any `json:"hourly"`
}

type WPoller struct {
	closech chan struct{}
	sender  Sender
}

func NewWPoller(sender Sender) *WPoller {
	return &WPoller{
		closech: make(chan struct{}),
		sender:  sender,
	}
}

func (wp *WPoller) Stop() {
	close(wp.closech)
}

type SmsSender struct {
	number string
}

func NewSmsSender(number string) *SmsSender {
	return &SmsSender{
		number: number,
	}
}

func (s *SmsSender) Send(data *WeatherData) error {
	fmt.Println("sending data to the number", s.number)
	return nil
}

func (wp *WPoller) start() {
	fmt.Println("starting WPoller")
	ticker := time.NewTicker(poller)
free:
	for {
		select {
		case <-ticker.C:
			data, err := getWeatherResults(28.44, 77.88)
			if err != nil {
				log.Fatal(err)
			}
			if err := wp.handleData(data); err != nil {
				log.Fatal(err)
			}

		case <-wp.closech:
			break free
		}
	}

	fmt.Println("wpoller stopped")
}

func (wp *WPoller) handleData(data *WeatherData) error {
	return wp.sender.Send(data)
}

func getWeatherResults(lat, long float64) (*WeatherData, error) {
	uri := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m", endpoint, lat, long)
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var data WeatherData
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func main() {
	smsSender := NewSmsSender("+917379745969")
	wpoller := NewWPoller(smsSender)
	go wpoller.start()

	time.Sleep(time.Minute * 1)

	wpoller.Stop()
}
