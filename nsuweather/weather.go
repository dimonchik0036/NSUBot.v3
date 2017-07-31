package nsuweather

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

func GetWeather() (weather string, err error) {
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

	if weather = string(reg.Find(body)); len(weather) < 2 {
		return "", errors.New("Weather not found")
	}

	return weather[1 : len(weather)-1], nil
}
