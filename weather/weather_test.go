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

// func TestMusicStyleByTemperature(t *testing.T) {
// 	// g := Goblin(t)
// 	var scenarios = []struct {
// 		in  float32
// 		out string
// 	}{
// 		{11, "rock"},
// 		{8, "classical"},
// 		{3, "party"},
// 	}

// 	client := weather.NewWeatherClient()

// 	for _, scenario := range scenarios {
// 		name := fmt.Sprintf("%f", scenario.in)
// 		t.Run(name, func(t *testing.T) {
// 			subject := client.MusicStyleByTemperature(scenario.in)
// 			if subject != scenario.out {
// 				// g.Assert(subject).Equal(scenario.out)
// 				t.Errorf("got %q, want %q", subject, scenario.out)
// 			}
// 		})
// 	}

// }

func TestMusicStyleByTemperature2(t *testing.T) {

	g := Goblin(t)
	// RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("MusicStyleByTemperature", func() {
		tests := []struct {
			it           string
			input        float32
			expectErr    string
			expectValues string
		}{
			{
				it:           "Temperature: 3 - when temperature is above 10 degrees, return classical music tracks",
				input:        -3,
				expectValues: "classical",
			},
			{
				it:           "Temperature: 9.999 - when temperature is above 10 degrees, return classical music tracks",
				input:        9.999,
				expectValues: "classical",
			},
			{
				it:           "Temperature: 11 - when between 10 and 14 degrees, return rock music tracks",
				input:        11,
				expectValues: "rock",
			},
			{
				it:           "Temperature: 14 - when between 10 and 14 degrees, return rock music tracks",
				input:        14,
				expectValues: "rock",
			},
			{
				it:           "Temperature: 14.4555 - when between 15 degrees and 30 degrees, return pop music tracks",
				input:        14.4555,
				expectValues: "rock",
			},
			{
				it:           "Temperature: 15 - when between 15 degrees and 30 degrees, return pop music tracks",
				input:        15,
				expectValues: "pop",
			},
			{
				it:           "Temperature: 30 - when between 15 degrees and 30 degrees, return pop music tracks",
				input:        30,
				expectValues: "pop",
			},
			{
				it:           "Temperature: 30.111 - when temperature is above 30 degrees, return party music tracks",
				input:        30.111,
				expectValues: "party",
			},
		}

		client := weather.NewWeatherClient()
		for _, test := range tests {
			test := test
			g.It(test.it, func() {
				subject := client.MusicStyleByTemperature(test.input)
				g.Assert(subject).Equal(test.expectValues)
			})
		}
	})
}
