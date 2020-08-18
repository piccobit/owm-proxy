package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	flag "github.com/spf13/pflag"
)

type owmData struct {
	Lat float64 `json:"lat,omitempty"`
	Lon float64 `json:"lon,omitempty"`
	Timezone string `json:"timezone,omitempty"`
	TimezoneOffset float32 `json:"timezone_offset,omitempty"`
	Current struct{
		Temp float64 `json:"temp,omitempty"`
		FeelsLike float64 `json:"feels_like,omitempty"`
		Pressure float64 `json:"pressure,omitempty"`
		Humidity int64 `json:"humidity"`
		DewPoint float64 `json:"dew_point,omitempty"`
		UVI float64 `json:"uvi,omitempty"`
		Clouds int64 `json:"clouds,omitempty"`
		Visibility int64 `json:"visibility,omitempty"`
		WindSpeed float64 `json:"wind_speed,omitempty"`
		WindDeg int64 `json:"wind_deg"`
	} `json:"current,omitempty"`
}

type owmSmallData struct {
	Temp float64 `json:"t"`
	FeelsLike float64 `json:"f"`
	Pressure float64 `json:"p"`
	Humidity int64 `json:"h"`
	DewPoint float64 `json:"d"`
	UVI float64 `json:"u"`
	Clouds int64 `json:"c"`
	Visibility int64 `json:"v"`
	WindSpeed float64 `json:"ws"`
	WindDeg int64 `json:"wd"`
}

var (
	version string = "dev"
	commit string = "none"
	date string = "unknown"
	debug bool
	port int64
	printVersion bool
)

func owm(ctx *gin.Context, lat string, lon string, units string, exclude string, appID string) {
	c := ctx.Copy()
	c = ctx

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?lat=%s&lon=%s&units=%s&exclude=%s&appid=%s", lat, lon, units, exclude, appID)

	resp, err := http.Get(url)

	if err != nil {
		log.Panicf("could not GET URL: %s", err.Error())
	}

	defer func() {
		err := resp.Body.Close()

		if err != nil {
			log.Panicf("could not close body reader: %s", err.Error())
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Panicf("could not read body: %s", err.Error())
	}

	var data owmData

	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Panicf("could not unmarshal body: %s", err.Error())
	}

	smallData := owmSmallData{
		data.Current.Temp,
		data.Current.FeelsLike,
		data.Current.Pressure,
		data.Current.Humidity,
		data.Current.DewPoint,
		data.Current.UVI,
		data.Current.Clouds,
		data.Current.Visibility,
		data.Current.WindSpeed,
		data.Current.WindDeg,
	}

	c.JSON(http.StatusOK, smallData)
}

func main() {
	flag.BoolVar(&printVersion, "version", false, "Print version")
	flag.BoolVar(&debug, "debug", false, "Enable debugging")
	flag.Int64Var(&port, "port", 8080, "Port to listen on")

	flag.Parse()

	if printVersion {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Commit:  %s\n", commit)
		fmt.Printf("Date:    %s\n", date)

		os.Exit(1)
	}

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	router.GET("/owm", func(c *gin.Context) {
		lat := c.Query("lat")
		lon := c.Query("lon")
		units := c.DefaultQuery("units", "metric")
		exclude := c.DefaultQuery("exclude", "minutely,hourly,daily")
		appID := c.Query("appid")

		log.Printf("appID: %s", appID)

		owm(c, lat, lon, units, exclude, appID)
	})

	listenAddr := fmt.Sprintf(":%d", port)

	err := router.Run(listenAddr)

	if err != nil {
		log.Panicf("could not start router: %s", err.Error())
	}
}
