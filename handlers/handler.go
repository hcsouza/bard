package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/hcsouza/bard/cache"
	"github.com/hcsouza/bard/injection"
	"github.com/hcsouza/bard/music"
	"github.com/hcsouza/bard/weather"
	"log"
	"net/http"
	"strconv"
)

func MusicByCityNameHandler(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)
	cityName := vars["name"]

	weatherClient := weather.NewWeatherClient()
	weather, err := weatherClient.WeatherByCityName(cityName)
	if err != nil {
		log.Println("Return Cached-Fallback")
	}

	genre := weatherClient.MusicStyleByTemperature(weather.Main.Temp)
	playlist, err := playlistByStyleAndCountry(genre, weather.Sys.Country)

	if err == nil {
		json.NewEncoder(w).Encode(playlist)
	} else {
		log.Println("Error on get Temperature: ", err)
		json.NewEncoder(w).Encode(struct {
			Music string `json:"music"`
		}{"none"})
	}
}

func MusicByCityCoordHandler(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(request)
	coords := weather.Coordinates{}

	value, err := strconv.ParseFloat(vars["lat"], 32)
	if err == nil {
		coords.Latitude = float32(value)
	} else {
		log.Println("Params wrong format: ", err)
	}
	value, err = strconv.ParseFloat(vars["lon"], 32)
	if err == nil {
		coords.Longitude = float32(value)
	} else {
		log.Println("Params wrong format: ", err)
	}

	weatherClient := weather.NewWeatherClient()
	weather, err := weatherClient.WeatherByCityCoord(coords)
	if err != nil {
		log.Println("Return Cached-Fallback")
	}
	genre := weatherClient.MusicStyleByTemperature(weather.Main.Temp)
	playlist, err := playlistByStyleAndCountry(genre, weather.Sys.Country)
	if err == nil {
		json.NewEncoder(w).Encode(playlist)
	} else {
		log.Println("Error on get Temperature: ", err)
		json.NewEncoder(w).Encode(struct {
			Music string `json:"music"`
		}{"none"})
	}
}

func playlistByStyleAndCountry(genre, country string) (music.Playlist, error) {
	cacheClient := injection.Get("CacheClient").(*cache.Client)
	result, err := cacheClient.TracksByCountryAndGenre(country, genre)
	if err == nil {
		return result, err
	}
	playlist, err := music.PlaylistByStyleAndCountry(genre, country)
	if err == nil {
		cacheClient.AddTracksByCountryAndGenre(country, genre, playlist)
	}
	return playlist, err
}
