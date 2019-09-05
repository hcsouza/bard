package main

import (
	"fmt"
	"github.com/gorilla/mux"
	. "github.com/hcsouza/bard/config"
	"github.com/hcsouza/bard/handlers"
	. "github.com/hcsouza/bard/logger"
	"github.com/hcsouza/bard/shield"
	"net/http"
)

func setupHandlers() *mux.Router {

	shield.StartMonitoring()

	r := mux.NewRouter()
	r.Path("/musics/city").
		Queries("name", "{name}").
		HandlerFunc(handlers.MusicByCityNameHandler).
		Name("ByName")

	r.HandleFunc("/musics/city", handlers.MusicByCityCoordHandler).
		Queries("lat", "{lat}", "lon", "{lon}").
		Name("ByCoord")

	return r
}

func main() {
	r := setupHandlers()
	port := fmt.Sprintf(":%s", Config.BardApi.Port)

	Logger.Message("Server is started!")
	Logger.Message(fmt.Sprintf("Listening on :%s", port))

	http.ListenAndServe(port, r)
}
