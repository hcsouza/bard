package weather

import (
	"encoding/json"
	"fmt"
	. "github.com/hcsouza/bard/config"
	"github.com/hcsouza/bard/shield"
	"log"
)

var (
	StylesByTemperature = map[string]string{
		"Above30": "party",
		"Above15": "pop",
		"Above10": "rock",
		"Below10": "classical",
	}
)

type weatherResult struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Main struct {
		Temp     float32 `json:"temp"`
		Pressure float32 `json:"pressure"`
		Humidity float32 `json:"humidity"`
		TempMin  float32 `json:"temp_min"`
		TempMax  float32 `json:"temp_max"`
	} `json:"main"`
	Sys struct {
		Country string `json:"country"`
	} `json:"sys"`
}

type Coordinates struct {
	Latitude  float32
	Longitude float32
}

type weatherClient struct{}

func NewWeatherClient() weatherClient {
	return weatherClient{}
}

func TemperatureByCityName(city string) (weather weatherResult, err error) {
	client := NewWeatherClient()
	weather, err = client.WeatherByCityName(city)
	return
}

func (client weatherClient) WeatherByCityName(city string) (weatherResult, error) {

	searchUrl := createUrlRequestByCityName(city)
	request := shield.CommandRequest{
		"TemperatureByCityName",
		searchUrl,
		"GET",
	}

	buffer, err := shield.ExecuteCommandWithCircuitBreaker(request)
	if err != nil {
		return weatherResult{}, err
	}

	return client.parseJsonToResult(buffer)
}

func (client weatherClient) WeatherByCityCoord(coords Coordinates) (weatherResult, error) {

	searchUrl := createUrlRequestByCoord(coords)
	request := shield.CommandRequest{
		"TemperatureByCityCoord",
		searchUrl,
		"GET",
	}

	buffer, err := shield.ExecuteCommandWithCircuitBreaker(request)
	if err != nil {
		return weatherResult{}, err
	}

	return client.parseJsonToResult(buffer)
}

func (client weatherClient) parseJsonToResult(jsonApi []byte) (result weatherResult, err error) {

	err = json.Unmarshal([]byte(jsonApi), &result)
	if err != nil {
		log.Println("Error on UnMarshall json Weather:  ", err)
		return
	}
	return result, err
}

func createUrlRequestByCoord(coords Coordinates) string {

	uriBase := Config.WeatherApi.Url
	apiKey := Config.WeatherApi.Appid

	search := fmt.Sprintf("%s?lat=%v", uriBase, coords.Latitude)
	search = fmt.Sprintf("%s&lon=%v&appid=%s", search, coords.Longitude, apiKey)
	search = fmt.Sprintf("%s&units=metric", search)

	return search
}

func createUrlRequestByCityName(cityName string) string {

	uriBase := Config.WeatherApi.Url
	apiKey := Config.WeatherApi.Appid

	search := fmt.Sprintf("%s?q=%s&appid=%s", uriBase, cityName, apiKey)
	search = fmt.Sprintf("%s&units=metric", search)

	return search
}

func (client weatherClient) MusicStyleByTemperature(temperature float32) string {

	switch {
	case temperature > 30:
		return StylesByTemperature["Above30"]
	case temperature > 15 && temperature <= 30:
		return StylesByTemperature["Above15"]
	case temperature > 10 && temperature <= 14.99:
		return StylesByTemperature["Above10"]
	default:
		return StylesByTemperature["Below10"]
	}
}
