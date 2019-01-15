package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WeatherResonse struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC      float32 `json:"temp_c"`
		FeelsLikeC float32 `json:"feelslike_c"`
		Humidity   int     `json:"humidity"`
		WindKPH    float32 `json:"wind_kph"`
		WindDir    string  `json:"wind_dir"`
	} `json:"current"`
}

type WeatherClient struct {
	URL string
	Key string
}

func (w WeatherClient) current() string {
	url := fmt.Sprintf("%s?key=%s&q=tel-aviv", w.URL, w.Key)
	response, _ := http.Get(url)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	weatherResponse := WeatherResonse{}
	json.Unmarshal(body, &weatherResponse)
	return fmt.Sprintf("It is %g degrees but feels like %g. A %s wind is blowing in a speed of %g KPH.\nHave a nice day", weatherResponse.Current.TempC, weatherResponse.Current.FeelsLikeC, weatherResponse.Current.WindDir, weatherResponse.Current.WindKPH)
}
