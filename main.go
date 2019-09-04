package main

import (
	"github.com/gorilla/mux"
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

	r.HandleFunc("/musics/coords", handlers.MusicByCityCoordHandler).
		Queries("lat", "{lat}", "lon", "{lon}").
		Name("ByCoord")

	return r
}

func main() {
	r := setupHandlers()

	Logger.Message("==> Main server is started")
	Logger.Message("listening on :8088")

	http.ListenAndServe(":8088", r)
}
