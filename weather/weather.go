package weather

import (
	"encoding/json"
	"fmt"
	"github.com/hcsouza/bard/shield"
	"io/ioutil"
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
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float32 `json:"temp_min"`
		TempMax  float32 `json:"temp_max"`
	} `json:"main"`
}

type weatherClient struct{}

func TemperatureByCityName(city string) string {

	client := NewWeatherClient()
	temp := client.TemperatureByCityName(city)

	return client.MusicStyleByTemperature(temp)
}

func NewWeatherClient() weatherClient {
	return weatherClient{}
}

func (client weatherClient) TemperatureByCityName(city string) (temperature float32) {
	var result weatherResult

	uriBase := "http://api.openweathermap.org/data/2.5/weather"
	apiKey := ""

	search := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", uriBase, city, apiKey)

	request := shield.CommandRequest{
		"TemperatureByCityName",
		search,
	}

	rp, err := shield.GetData(request)
	buffBytes, err := ioutil.ReadAll(rp.Body)
	if err != nil {
		fmt.Println("err: ", err)
	}
	_ = json.Unmarshal([]byte(buffBytes), &result)
	temperature = result.Main.Temp

	return
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
