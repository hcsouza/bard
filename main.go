package main

import (
	"fmt"
	"github.com/hcsouza/bard/weather"
)

func main() {
	style := weather.TemperatureByCityName("caçapava")
	fmt.Println("Music Style: ", style)
}
