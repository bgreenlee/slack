package handlers

/**
 *   /weather [city]  -  Return the current weather for the given city
 */

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bradfitz/latlong"

	"github.com/bgreenlee/slack/config"
)

/**
 * weatherHandler reports the weather for a given location using OpenWeatherMap
 */
func Weather(resp http.ResponseWriter, req *http.Request, config config.Config) {
	isslackrequest := req.FormValue("token") == config.SlackToken

	// Slack sends the text after the command in the "text" query param
	location := req.FormValue("text")

	// make our weather API call
	params := url.Values{}
	params.Add("q", location)
	params.Add("units", "imperial") // suck it, rest of the world
	res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?" + params.Encode())

	// read the response and parse the JSON
	decoder := json.NewDecoder(res.Body)
	var weatherData struct {
		Coord *struct {
			Lon float64
			Lat float64
		}
		Sys *struct {
			SunriseEpoch int64 `json:"sunrise"`
			SunsetEpoch  int64 `json:"sunset"`
		}
		Weather []*struct {
			Description string
			Icon        string
		}
		Main *struct {
			Temp     float64
			TempHigh float64 `json:"temp_max"`
			TempLow  float64 `json:"temp_min"`
		}
		Location string      `json:"name"`
		Message  string      // only returned in the case of an error
		Code     interface{} `json:"cod"` // this can be either an int or a string, ugh
	}
	err = decoder.Decode(&weatherData)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(resp, "Oops, got an error: %s", err)
		return
	}

	if weatherData.Main == nil {
		if weatherData.Code.(string) == "404" {
			fmt.Fprintf(resp, "Sorry, I couldn't find anything for '%s'", location)
		} else {
			fmt.Fprintf(resp, "Oops, got an error: %s", weatherData.Message)
		}
		return
	}

	// build the conditions string
	conditions := []string{}
	conditionsStr := ""
	for _, conditionData := range weatherData.Weather {
		conditions = append(conditions, conditionData.Description)
	}
	if len(conditions) > 2 {
		conditionsStr = strings.Join(conditions[:len(conditions)-1], ", ") + ", and " + conditions[len(conditions)-1]
	} else {
		conditionsStr = strings.Join(conditions, " and ")
	}

	// get sunrise/sunset in the local time zone
	timezone := latlong.LookupZoneName(weatherData.Coord.Lat, weatherData.Coord.Lon)
	tzLocation, _ := time.LoadLocation(timezone)
	sunrise := time.Unix(weatherData.Sys.SunriseEpoch, 0).In(tzLocation)
	sunset := time.Unix(weatherData.Sys.SunsetEpoch, 0).In(tzLocation)

	// send our response
	responseStr := fmt.Sprintf("Currently %d and %s. The high today is %d and the low is %d.\nSunrise is at %s and sunset is at %s",
		int(weatherData.Main.Temp), conditionsStr,
		int(weatherData.Main.TempHigh), int(weatherData.Main.TempLow),
		sunrise.Format("3:04am"), sunset.Format("3:04pm"))

	// send the response back to Slack
	if isslackrequest {
		postResponse(config.ResponseURL, responseStr, "Weather for "+weatherData.Location, "http://openweathermap.org/img/w/"+weatherData.Weather[0].Icon+".png")
	} else {
		fmt.Fprint(resp, responseStr)
	}
}
