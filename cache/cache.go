package cache

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	. "github.com/hcsouza/bard/config"
	"github.com/hcsouza/bard/music"
	"log"
)

type Client struct {
	memcacheClient *memcache.Client
}

func NewCacheClient() *Client {
	mc := memcache.New(Config.Memcache.Host)
	return &Client{mc}
}

func (client *Client) AddTracksByCountryAndGenre(country, genre string, playlist music.Playlist) error {

	bytes, err := json.Marshal(&playlist)
	if err != nil {
		log.Println("Error on serialize tracks: ", err)
	}
	itemKey := fmt.Sprintf("%s:%s", country, genre)
	serialized := &memcache.Item{Key: itemKey, Value: []byte(bytes)}

	err = client.memcacheClient.Set(serialized)
	if err != nil {
		log.Println("Error on Set serialized to cache: ", err)
		return err
	}
	msg := fmt.Sprintf("Tracks for genre: %s on country: %s, adde to cache.", genre, country)
	log.Println(msg)
	return err
}

func (client *Client) TracksByCountryAndGenre(country, genre string) (music.Playlist, error) {
	var result music.Playlist
	searchKey := fmt.Sprintf("%s:%s", country, genre)

	it, err := client.memcacheClient.Get(searchKey)
	if err != nil {
		log.Println("Error on Get from memcache:, ", err)
		return result, err
	}

	err = json.Unmarshal([]byte(it.Value), &result)
	if err != nil {
		log.Println("Error on UnMarshall from memcache:, ", err)
		return result, err
	}
	msg := fmt.Sprintf("Tracks for genre: %s on country: %s, founded on cache.", genre, country)
	log.Println(msg)
	return result, err
}
