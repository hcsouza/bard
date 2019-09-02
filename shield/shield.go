package shield

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"gopkg.in/eapache/go-resiliency.v1/retrier"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type CommandRequest struct {
	Name   string
	Url    string
	Method string
}

func init() {
	setupCommands()
}

func setupCommands() {
	hystrix.ConfigureCommand("TemperatureByCityName", hystrix.CommandConfig{
		Timeout:                5000,
		MaxConcurrentRequests:  300,
		RequestVolumeThreshold: 3,
		SleepWindow:            1000,
		ErrorPercentThreshold:  10,
	})

	hystrix.ConfigureCommand("TemperatureByCityCoord", hystrix.CommandConfig{
		Timeout:                50000,
		MaxConcurrentRequests:  300,
		RequestVolumeThreshold: 1,
		SleepWindow:            1000,
		ErrorPercentThreshold:  10,
	})
}

func StartMonitoring() {
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)
}

func ExecuteCommandWithCircuitBreaker(commandRequest CommandRequest) ([]byte, error) {

	chSuccess := make(chan []byte, 1)
	errors := hystrix.Go(commandRequest.Name,
		func() error {
			request, _ := http.NewRequest(commandRequest.Method, commandRequest.Url, nil)
			return DoCallRequestWithRetries(request, chSuccess)
		},
		func(err error) error {
			log.Println(fmt.Sprintf("Fallback for %s, with error: %s", commandRequest.Name, err.Error()))
			return err
		})

	select {
	case out := <-chSuccess:
		log.Println("Successful call on", commandRequest.Name)
		return out, nil

	case err := <-errors:
		return nil, err
	}
}

func DoCallRequestWithRetries(commandRequest *http.Request, chSuccess chan []byte) error {

	r := retrier.New(retrier.ConstantBackoff(3, 100*time.Millisecond), nil)
	client := &http.Client{}

	err := r.Run(func() error {
		response, err := client.Do(commandRequest)
		if err == nil && response.StatusCode == 200 {
			buffer, err := ioutil.ReadAll(response.Body)
			if err == nil {
				chSuccess <- buffer
				return nil
			}
			return err
		} else if err == nil {
			err = fmt.Errorf("Status was %v", response.StatusCode)
		}
		return err
	})
	return err
}
