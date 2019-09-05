package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
	"github.com/hcsouza/bard/cache"
	"github.com/hcsouza/bard/injection"
	. "github.com/hcsouza/bard/logger"
	"github.com/hcsouza/bard/music"
	"github.com/hcsouza/bard/weather"
	"net/http"
	"strconv"
)

var (
	ErrFallbackTracks = errors.New("Fallback: can't get from cache or music api")
)

func MusicByCityNameHandler(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)
	cityName := vars["name"]

	weatherClient := weather.NewWeatherClient()
	weather, err := weatherClient.WeatherByCityName(cityName)
	if err != nil {
		Logger.Warn("Return Cached-Fallback")
		playlist, _ := GetFallBackPlayList()
		json.NewEncoder(w).Encode(playlist)
		return
	}

	genre := weatherClient.MusicStyleByTemperature(weather.Main.Temp)
	playlist, err := playlistByStyleAndCountry(genre, weather.Sys.Country)

	if err == ErrFallbackTracks {
		Logger.Error(err.Error())
	}
	json.NewEncoder(w).Encode(playlist)
}

func MusicByCityCoordHandler(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)
	coords := weather.Coordinates{}

	value, err := strconv.ParseFloat(vars["lat"], 32)
	if err == nil {
		coords.Latitude = float32(value)
	} else {
		Logger.Error(fmt.Sprintf("Params wrong format: %s", err))
		badrequest, _ := json.Marshal(struct {
			Message string `json:"message"`
		}{"Params wrong format."})

		http.Error(w, string(badrequest), http.StatusBadRequest)
		return
	}
	value, err = strconv.ParseFloat(vars["lon"], 32)
	if err == nil {
		coords.Longitude = float32(value)
	} else {
		Logger.Error(fmt.Sprintf("Params wrong format: %s", err))
		badrequest, _ := json.Marshal(struct {
			Message string `json:"message"`
		}{"Params wrong format."})

		http.Error(w, string(badrequest), http.StatusBadRequest)
		return
	}

	weatherClient := weather.NewWeatherClient()
	weather, err := weatherClient.WeatherByCityCoord(coords)
	if err != nil {
		Logger.Warn("Return Cached-Fallback")
		playlist, _ := GetFallBackPlayList()
		json.NewEncoder(w).Encode(playlist)
		return
	}
	genre := weatherClient.MusicStyleByTemperature(weather.Main.Temp)
	playlist, err := playlistByStyleAndCountry(genre, weather.Sys.Country)

	if err == ErrFallbackTracks {
		Logger.Error(err.Error())
	}
	json.NewEncoder(w).Encode(playlist)
}

func playlistByStyleAndCountry(genre, country string) (music.Playlist, error) {
	cacheClient := injection.Get("CacheClient").(*cache.Client)
	playlist, err := cacheClient.TracksByCountryAndGenre(country, genre)

	if err == nil {
		return playlist, err
	}
	if err == memcache.ErrCacheMiss {
		playlist, err = music.PlaylistByStyleAndCountry(genre, country)
		if err == nil {
			cacheClient.AddTracksByCountryAndGenre(country, genre, playlist)
		}
		return playlist, err
	}
	Logger.Info("Returning FallBackList")
	return GetFallBackPlayList()
}

func GetFallBackPlayList() (music.Playlist, error) {
	playlist, err := music.PlaylistByStyleAndCountry("rock", "us")
	if err != nil {
		Logger.Error(fmt.Sprintf("Error on get FallbackList on music api: %s", err))
		musicOne := music.Music{"Patience", "Guns N' Roses"}
		return music.Playlist{[]music.Music{musicOne}}, ErrFallbackTracks
	}
	return playlist, err
}
