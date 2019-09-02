package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
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
	temperature, err := weatherClient.TemperatureByCityName(cityName)
	if err != nil {
		log.Println("Return Cached-Fallback")
	}
	styleName := weatherClient.MusicStyleByTemperature(temperature)

	//Search from SpotifyApp
	//#Success
	// return Json, with track list
	//#fail
	// return cached-fallback

	//Retorno do Handler
	data := struct {
		Music string `json:"music"`
	}{
		styleName,
	}

	if err == nil {
		json.NewEncoder(w).Encode(data)
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
	temperature, err := weatherClient.TemperatureByCityCoord(coords)
	if err != nil {
		log.Println("Return Cached-Fallback")
	}
	styleName := weatherClient.MusicStyleByTemperature(temperature)

	//Search from SpotifyApp
	//#Success
	// return Json, with track list
	//#fail
	// return cached-fallback

	//Retorno do Handler
	data := struct {
		Music string `json:"music"`
	}{
		styleName,
	}

	if err == nil {
		json.NewEncoder(w).Encode(data)
	} else {
		log.Println("Error on get Temperature: ", err)
		json.NewEncoder(w).Encode(struct {
			Music string `json:"music"`
		}{"none"})
	}
}
