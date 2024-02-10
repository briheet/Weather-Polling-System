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
	senders []Sender
}

func NewWPoller(senders ...Sender) *WPoller {
	return &WPoller{
		closech: make(chan struct{}),
		senders: senders,
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

type emailSender struct {
	email string
}

func NewEmailSender(email string) *emailSender {
	return &emailSender{
		email: email,
	}
}

func (s *SmsSender) Send(data *WeatherData) error {
	fmt.Println("sending data to the number", s.number)
	return nil
}

func (em *emailSender) Send(data *WeatherData) error {
	fmt.Println("sending data to the email", em.email)
	return nil
}

//func (wp *WPoller) start() {
//	fmt.Println("starting WPoller")
//	ticker := time.NewTicker(poller)
//free:
//	for {
//		select {
//		case <-ticker.C:
//			data, err := getWeatherResults(28.44, 77.88)
//			if err != nil {
//				log.Fatal(err)
//			}
//			if err := wp.handleData(data); err != nil {
//				log.Fatal(err)
//			}
//
//		case <-wp.closech:
//			break free
//		}
//	}
//
//	fmt.Println("wpoller stopped")
//}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	lat := 28.44
	long := 77.88

	data, err := getWeatherResults(lat, long)
	if err != nil {
		http.Error(w, "Failed to get weather data", http.StatusInternalServerError)
		return
	}

	currentTime := time.Now().Format("2006-01-02T15:00") // Get current time in the required format

	// Find index of the current time in the data
	index := -1
	for i, t := range data.Hourly["time"].([]interface{}) {
		if t.(string) == currentTime {
			index = i
			break
		}
	}

	if index != -1 {
		temp := data.Hourly["temperature_2m"].([]interface{})[index].(float64)
		fmt.Printf("Current temperature at %s is %.2f\n", currentTime, temp)

		// Create a new WeatherData instance with only current time and temperature
		currentWeather := &WeatherData{
			Elevation: data.Elevation,
			Hourly: map[string]any{
				"time":        currentTime,
				"temperature": temp,
			},
		}

		// Encode and send the current weather data
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentWeather)
	} else {
		fmt.Printf("Current time %s not found in weather data\n", currentTime)
		http.Error(w, "Current time not found in weather data", http.StatusNotFound)
	}
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
	//	smsSender := NewSmsSender("+917379745969")
	//	emailSender := NewEmailSender("briheetyadav@gmail.com")
	//	wpoller := NewWPoller(smsSender, emailSender)
	//	wpoller.start()

	fmt.Println("starting port at :8080")
	http.HandleFunc("/weather", weatherHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
