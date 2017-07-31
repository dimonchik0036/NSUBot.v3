package nsuweather

import (
	"regexp"
	"testing"
)

func TestGetWeather(t *testing.T) {
	w, err := GetWeather()
	if err != nil {
		t.Fatal(err)
	}

	reg, err := regexp.Compile("НГУ")
	if err != nil {
		t.Fatal(err)
	}

	if check := reg.FindString(w); check == "" {
		t.Fatal("Не найден НГУ")
	}
}
