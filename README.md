# Bard

Bard is a rest-api to list tracks using location and weather.

Basically the consumer can get weather information by city name or coordinates, and then use the weather result and country to filter tracks.

This api consume two external services https://openweathermap.org for weather and  https://developer.spotify.com for music tracks. For deal with this dependecies the api was builded using circuit-breaker pattern, retries and cache.

### Resources

#### Tracks by City Name

**GET** http://$endpoint_api:8088/musics/city?name=london

Response: 200

```json
{
   "musics":[
      {
         "name":"Señorita",
         "artist":"Shawn Mendes"
      },
      {
         "name":"Take Me Back to London (feat. Stormzy)",
         "artist":"Ed Sheeran"
      },
      {
         "name":"Beautiful People (feat. Khalid)",
         "artist":"Ed Sheeran"
      },
      {
         "name":"How Do You Sleep?",
         "artist":"Sam Smith"
      },
      {
         "name":"Higher Love",
         "artist":"Kygo"
      }
   ]
}
 ``` 

#### Tracks by City Coords

**GET** http://$endpoint_api:8088/musics/city?lat=51.51&lon=-0.13

Response: 200

```json
{
   "musics":[
      {
         "name":"Señorita",
         "artist":"Shawn Mendes"
      },
      {
         "name":"Take Me Back to London (feat. Stormzy)",
         "artist":"Ed Sheeran"
      },
      {
         "name":"Beautiful People (feat. Khalid)",
         "artist":"Ed Sheeran"
      },
      {
         "name":"How Do You Sleep?",
         "artist":"Sam Smith"
      },
      {
         "name":"Higher Love",
         "artist":"Kygo"
      }
   ]
}
 ``` 


### Tech

Bard api was built using [Golang](https://golang.org/) and some tools:
* [Gorilla Mux](http://www.gorillatoolkit.org/pkg/mux) - Web toolkit and http router
 * [hystrix-go](https://github.com/afex/hystrix-go) - Circuit-breaker tool
 * [retrier](https://godoc.org/gopkg.in/eapache/go-resiliency.v1/retrier) - Retrier package
 * [zmb3/spotify](https://github.com/zmb3/spotify) - Spotify Client
 * [gomemcache](https://github.com/bradfitz/gomemcache) - Memcache Client
 * [gologger](https://github.com/sadlil/gologger) - Logger library
 * [minimock](https://github.com/gojuno/minimock) - Mock library
 * [goblin](https://github.com/franela/goblin) - BDD testing framework
 * [sarulabs/di](https://github.com/sarulabs/di) - Dependencie injection container
 * [viper](https://github.com/spf13/viper) - Config solution tool


## Run

The app can be started using docker or kubernetes.

* ***Docker***
Using compose file of project with `web` parameter.

```sh
$ docker-compose up web
```
or using the setup docker shell script:

```sh
./docker-setup.sh
```
To see all access points you can run another sh script:
```sh
./docker-access-points.sh
```
This script will print all access points of project:
1. RestApi resource
2. RestApi Stream
3. Hystrix Dashboard
---
* **Kubernetes**

With an environemt with kubernetes ( a suggestion is [Minikube](https://github.com/kubernetes/minikube) ) you need to run the setup script:
```sh
./k8s-setup.sh
```
To see all access points you can run another sh script:
```sh
./k8s-access-points.sh  ClusterName
```
This script will print all access points of project:
1. RestApi resource
2. RestApi Stream
3. Hystrix Dashboard


### Development

To run in development for debug or improvement you can use another docker parameter:
```sh
docker-compose up development
```
And connect to container:
```sh
docker-compose exec development /bin/bash
```

### Test

In a development container just run:

```sh
ENV=test go test ./...  -v
```
