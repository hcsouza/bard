package weather_test

import (
	"encoding/json"
	"fmt"
	. "github.com/franela/goblin"
	"github.com/hcsouza/bard/logger"
	"github.com/hcsouza/bard/weather"
	"github.com/nbio/st"
	"github.com/sadlil/gologger"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

func init() {
	newLogger := gologger.GetLogger(gologger.FILE, "/dev/null")
	logger.SetLogger(newLogger)
}

func TestWeatherByCityName(t *testing.T) {
	defer gock.Off()
	g := Goblin(t)

	g.Describe("WeatherByCityName", func() {
		g.It("with a valid city name should success", func() {

			var expected weather.WeatherResult
			bodyResponse := `{"coord":{"lon":-0.13,"lat":51.51},"weather":[{"id":803,"main":"Clouds","description":"broken clouds","icon":"04d"}],"base":"stations","main":{"temp":289.58,"pressure":1021,"humidity":59,"temp_min":288.15,"temp_max":291.15},"visibility":10000,"wind":{"speed":7.2,"deg":240},"clouds":{"all":75},"dt":1567771289,"sys":{"type":1,"id":1414,"message":0.0141,"country":"GB","sunrise":1567747236,"sunset":1567795072},"timezone":3600,"id":2643743,"name":"cacapava","cod":200}`

			url := "http://api.openweathermap.org/data/2.5/weather"
			gock.New(url).
				MatchParam("q", "cacapava").
				MatchParam("appid", "1234").
				MatchParam("units", "metric").
				Reply(200).
				BodyString(bodyResponse)

			subject := weather.NewWeatherClient()
			result, _ := subject.WeatherByCityName("cacapava")

			_ = json.Unmarshal([]byte(bodyResponse), &expected)
			g.Assert(result).Equal(expected)
		})
	})
	st.Expect(t, gock.IsDone(), true)
}

func TestWeatherByCityCoords(t *testing.T) {
	defer gock.Off()
	g := Goblin(t)

	g.Describe("WeatherByCityCoord", func() {
		g.It("with a valid city coords should success", func() {

			var expected weather.WeatherResult
			bodyResponse := `{"coord":{"lon":-0.13,"lat":51.51},"weather":[{"id":803,"main":"Clouds","description":"broken clouds","icon":"04d"}],"base":"stations","main":{"temp":289.58,"pressure":1021,"humidity":59,"temp_min":288.15,"temp_max":291.15},"visibility":10000,"wind":{"speed":7.2,"deg":240},"clouds":{"all":75},"dt":1567771289,"sys":{"type":1,"id":1414,"message":0.0141,"country":"GB","sunrise":1567747236,"sunset":1567795072},"timezone":3600,"id":2643743,"name":"cacapava","cod":200}`
			url := "http://api.openweathermap.org/data/2.5/weather"

			gock.New(url).
				MatchParam("lat", "-0.13").
				MatchParam("lon", "51.51").
				MatchParam("appid", "1234").
				MatchParam("units", "metric").
				Reply(200).
				BodyString(bodyResponse)

			subject := weather.NewWeatherClient()
			params := weather.Coordinates{-0.13, 51.51}

			result, _ := subject.WeatherByCityCoord(params)

			_ = json.Unmarshal([]byte(bodyResponse), &expected)
			g.Assert(result).Equal(expected)
		})
	})
	st.Expect(t, gock.IsDone(), true)
}

func TestCreateUrlRequestByCityName(t *testing.T) {

	g := Goblin(t)
	g.Describe("CreateUrlRequestByCityName", func() {
		g.It("make a valid url to request weather", func() {

			uriBase := "http://api.openweathermap.org/data/2.5/weather"
			expected := fmt.Sprintf("%s?q=cacapava&appid=1234&units=metric", uriBase)

			result := weather.CreateUrlRequestByCityName("cacapava")
			g.Assert(result).Equal(expected)
		})
	})
}

func TestCreateUrlRequestByCityCoords(t *testing.T) {

	g := Goblin(t)
	g.Describe("CreateUrlRequestByCoord", func() {
		g.It("make a valid url to request weather", func() {

			uriBase := "http://api.openweathermap.org/data/2.5/weather"
			expected := fmt.Sprintf("%s?lat=-0.13&lon=51.51&appid=1234&units=metric", uriBase)
			coords := weather.Coordinates{-0.13, 51.51}

			result := weather.CreateUrlRequestByCoord(coords)
			g.Assert(result).Equal(expected)
		})
	})
}
