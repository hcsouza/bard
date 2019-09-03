package music

import (
	"context"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	. "github.com/hcsouza/bard/config"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"log"
)

type Music struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
}

type Playlist struct {
	Musics []Music `json:"musics"`
}

func PlaylistByStyleAndCountry(musicStyle string, country string) (Playlist, error) {

	config := &clientcredentials.Config{
		ClientID:     Config.SpotifyApi.ClientId,
		ClientSecret: Config.SpotifyApi.SecretKey,
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Println("Couldn't get app music token: ", err)
	}

	client := spotify.Authenticator{}.NewClient(token)

	limit := 5
	opts := &spotify.Options{Country: &country, Limit: &limit}
	query := fmt.Sprintf("genre:%s", musicStyle)

	chSuccess := make(chan *spotify.SearchResult, 1)
	errors := hystrix.Go("PlaylistByCountry",
		func() error {
			results, err := client.SearchOpt(query, spotify.SearchTypeTrack, opts)
			if err != nil {
				log.Println("Error on Search Tracks: ", err)
				return err
			}

			chSuccess <- results
			return nil
		},
		func(err error) error {
			log.Println(fmt.Sprintf("Fallback for %s, with error: %s", "PlaylistByCountry", err.Error()))
			return err
		})

	select {
	case out := <-chSuccess:
		log.Println("Successful call on PlaylistByCountry")
		return parseResultToPlayList(out)
	case err := <-errors:
		return Playlist{}, err
	}
}

func parseResultToPlayList(results *spotify.SearchResult) (Playlist, error) {
	var playList Playlist
	var err error

	if results.Tracks != nil {
		for _, item := range results.Tracks.Tracks {
			music := Music{item.Name, item.Artists[0].Name}
			playList.Musics = append(playList.Musics, music)
		}
	} else {
		err = errors.New("Track result is empty")
	}
	return playList, err
}
