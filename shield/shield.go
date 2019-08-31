package shield

import (
	"github.com/afex/hystrix-go/hystrix"
	"net/http"
)

type CommandRequest struct {
	Name string
	Url  string
}

func init() {
	setupCommands()
}

func setupCommands() {
	hystrix.ConfigureCommand("TemperatureByCityName", hystrix.CommandConfig{
		Timeout:                50000,
		MaxConcurrentRequests:  300,
		RequestVolumeThreshold: 10,
		SleepWindow:            1000,
		ErrorPercentThreshold:  50,
	})

	hystrix.ConfigureCommand("TemperatureByCityCoord", hystrix.CommandConfig{
		Timeout:                50000,
		MaxConcurrentRequests:  300,
		RequestVolumeThreshold: 10,
		SleepWindow:            1000,
		ErrorPercentThreshold:  50,
	})
}

func GetData(request CommandRequest) (*http.Response, error) {
	resp, err := http.Get(request.Url)
	if err != nil {
		return resp, err
	}
	return resp, err
}
