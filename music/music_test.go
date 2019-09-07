package music_test

import (
	// "errors"
	"github.com/bradfitz/gomemcache/memcache"
	. "github.com/franela/goblin"
	"github.com/gojuno/minimock"
	"github.com/hcsouza/bard/cache"
	"github.com/hcsouza/bard/injection"
	"github.com/hcsouza/bard/logger"
	"github.com/hcsouza/bard/mock"
	"github.com/hcsouza/bard/music"
	"github.com/nbio/st"
	"github.com/sadlil/gologger"
	"github.com/sarulabs/di"
	"github.com/zmb3/spotify"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

func init() {
	newLogger := gologger.GetLogger(gologger.FILE, "/dev/null")
	logger.SetLogger(newLogger)
}

func TestWeatherByCityName(t *testing.T) {
	defer gock.Off()
	gock.Observe(gock.DumpRequest)
	mc := minimock.NewController(t)
	defer mc.Finish()
	g := Goblin(t)

	g.Describe("WeatherByCityName", func() {
		g.BeforeEach(func() {
			InjectServicesTesting(t, mc)
		})
		g.It("with a valid city name should success", func() {

			var expected music.Playlist
			item := music.Music{"Gravity", "Led Zeppelin"}
			expected.Musics = append(expected.Musics, item)

			service := injection.Get("MusicClientSearcher").(music.MusicService)
			subject := music.NewMusicClient(service)

			res, err := subject.PlaylistByStyleAndCountry("pop", "br")
			g.Assert(err).Equal(nil)
			g.Assert(res).Equal(expected)
		})
	})
	st.Expect(t, gock.IsDone(), true)
}

func InjectServicesTesting(t *testing.T, mc minimock.MockController) {
	services := []di.Def{
		{
			Name:  "CacheClient",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				fake := ExpectationsCacheMusicServiceMock(mc)
				return fake, nil
			},
		},
		{
			Name:  "MusicClientSearcher",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				fake := ExpectationsMusicServiceMock(mc)
				return fake, nil
			},
		},
	}

	container := injection.CreateContainer(services)
	injection.SetContainerApp(container)

}

func ExpectationsCacheMusicServiceMock(mc minimock.MockController) cache.CacheMusicService {

	result := music.Playlist{}
	musicServiceMock := mock.NewCacheMusicServiceMock(mc).
		TracksByCountryAndGenreMock.
		When("pop", "br").
		Then(result, memcache.ErrCacheMiss)

	return musicServiceMock
}

func ExpectationsMusicServiceMock(mc minimock.MockController) music.MusicService {
	country := "br"
	limit := 5
	opts := &spotify.Options{Country: &country, Limit: &limit}

	item := spotify.FullTrack{}
	item.Name = "Gravity"
	item.Artists = []spotify.SimpleArtist{spotify.SimpleArtist{Name: "Led Zeppelin"}}
	tracks := []spotify.FullTrack{item}
	page := spotify.FullTrackPage{Tracks: tracks}
	result := &spotify.SearchResult{
		Tracks: &page,
	}

	musicServiceMock := mock.NewMusicServiceMock(mc).
		SearchOptMock.
		When("genre:pop", 8, opts).
		Then(result, nil)

	return musicServiceMock
}
