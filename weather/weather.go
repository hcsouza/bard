package weather

import (
	"encoding/json"
	"fmt"
	"github.com/hcsouza/bard/shield"
	"log"
	"math/rand"
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
}

type Coordinates struct {
	Latitude  float32
	Longitude float32
}

type weatherClient struct{}

func NewWeatherClient() weatherClient {
	return weatherClient{}
}

func TemperatureByCityName(city string) (temperature float32, err error) {
	client := NewWeatherClient()
	temperature, err = client.TemperatureByCityName(city)
	return
}

func (client weatherClient) TemperatureByCityName(city string) (temperature float32, err error) {
	var result weatherResult
	var uriBase string

	chaosMonkey := rand.Intn(10)
	if chaosMonkey > 1 {
		uriBase = "http://apiaaa.openweathermap.org/data/2.5/weather"
	} else {
		uriBase = "http://api.openweathermap.org/data/2.5/weather"
	}
	log.Println(uriBase)

	apiKey := "69075f27ec95ce1dcd970bfc4eb5233f"
	search := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", uriBase, city, apiKey)

	request := shield.CommandRequest{
		"TemperatureByCityName",
		search,
		"GET",
	}

	buffer, err := shield.ExecuteCommandWithCircuitBreaker(request)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal([]byte(buffer), &result)
	if err != nil {
		log.Println("Error on UnMarshall json Weather:  ", err)
	}
	temperature = result.Main.Temp

	return temperature, err
}

func (client weatherClient) TemperatureByCityCoord(coords Coordinates) (temperature float32, err error) {
	var result weatherResult
	var uriBase string

	chaosMonkey := rand.Intn(10)
	if chaosMonkey > -2 {
		uriBase = "http://apiaaa.openweathermap.org/data/2.5/weather"
	} else {
		uriBase = "http://api.openweathermap.org/data/2.5/weather"
	}
	log.Println(uriBase)

	apiKey := "69075f27ec95ce1dcd970bfc4eb5233f"
	search := fmt.Sprintf("%s?lat=%v&lon=%v&appid=%s&units=metric", uriBase, coords.Latitude, coords.Longitude, apiKey)

	request := shield.CommandRequest{
		"TemperatureByCityCoord",
		search,
		"GET",
	}

	buffer, err := shield.ExecuteCommandWithCircuitBreaker(request)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal([]byte(buffer), &result)
	if err != nil {
		log.Println("Error on UnMarshall json Weather:  ", err)
	}
	temperature = result.Main.Temp

	return temperature, err

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
