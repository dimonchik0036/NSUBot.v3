package nsuweather

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

type Weather struct {
	Mux     sync.Mutex
	Weather string
}

func GetWeather() (string, error) {
	res, err := http.Get("http://weather.nsu.ru/loadata.php")
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", errors.New("Status error: " + res.Status)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	reg, err := regexp.Compile("'Температура около .*?'")
	if err != nil {
		return "", err
	}

	currentWeather := reg.Find(body)
	if len(currentWeather) < 2 {
		return "", errors.New("Weather not found")
	}

	return string(currentWeather), nil
}

func NewWeather() (weather Weather) {
	weather.Weather, _ = GetWeather()
	return
}

func (weather *Weather) Update() {
	currentWeather, err := GetWeather()
	if err != nil {
		return
	}

	weather.Mux.Lock()
	weather.Weather = currentWeather
	weather.Mux.Unlock()
}
