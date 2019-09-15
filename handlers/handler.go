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
	weat, err := weatherClient.WeatherByCityName(cityName)

	if err != nil {
		Logger.Warn("Return Cached-Fallback")
		json.NewEncoder(w).Encode(GetFallBackPlayList())
		return
	}

	playlist, err := playlistByStyleAndCountry(
		weatherClient.MusicStyleByTemperature(weat.Main.Temp),
		weat.Sys.Country)

	json.NewEncoder(w).Encode(playlist)
}

func MusicByCityCoordHandler(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)

	coords, err := checkCoords(vars)
	if err != nil {
		BadRequestCoords(w)
		return
	}

	weatherClient := weather.NewWeatherClient()
	weather, err := weatherClient.WeatherByCityCoord(coords)
	if err != nil {
		Logger.Warn("Return Cached-Fallback")
		json.NewEncoder(w).Encode(GetFallBackPlayList())
		return
	}

	playlist, err := playlistByStyleAndCountry(
		weatherClient.MusicStyleByTemperature(weather.Main.Temp),
		weather.Sys.Country)

	json.NewEncoder(w).Encode(playlist)
}

func DescribeResources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "describe.json")
}

func checkCoords(vars map[string]string) (weather.Coordinates, error) {
	coords := weather.Coordinates{}
	value, err := strconv.ParseFloat(vars["lat"], 32)
	if err == nil {
		coords.Latitude = float32(value)
	} else {
		return coords, errors.New(fmt.Sprintf("Params wrong format: %s", err))
	}

	value, err = strconv.ParseFloat(vars["lon"], 32)
	if err == nil {
		coords.Longitude = float32(value)
	} else {
		return coords, errors.New(fmt.Sprintf("Params wrong format: %s", err))
	}
	return coords, err
}

func BadRequestCoords(w http.ResponseWriter) {
	badrequest, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{"Params wrong format."})

	http.Error(w, string(badrequest), http.StatusBadRequest)
	return
}

func playlistByStyleAndCountry(genre, country string) (music.Playlist, error) {

	cacheClient := injection.Get("CacheClient").(*cache.Client)
	playlist, err := cacheClient.TracksByCountryAndGenre(country, genre)

	service := injection.Get("MusicClientSearcher").(music.MusicService)
	musicClient := music.NewMusicClient(service)

	if err == nil {
		return playlist, err
	} else {
		if err == memcache.ErrCacheMiss {
			playlist, err = musicClient.PlaylistByStyleAndCountry(genre, country)
			if err == nil {
				_ = cacheClient.AddTracksByCountryAndGenre(country, genre, playlist)
				return playlist, err
			}
		}
	}
	Logger.Info("Can't get track from cache and apis - Returning FallBackList")
	return GetFallBackPlayList(), err
}

func GetFallBackPlayList() music.Playlist {

	service := injection.Get("MusicClientSearcher").(music.MusicService)
	musicClient := music.NewMusicClient(service)
	playlist, err := musicClient.PlaylistByStyleAndCountry("rock", "us")

	if err != nil {
		Logger.Error(fmt.Sprintf("Error on get FallbackList on music api: %s", err))
		musicOne := music.Music{"Patience", "Guns N' Roses"}
		return music.Playlist{[]music.Music{musicOne}}
	}
	return playlist
}
