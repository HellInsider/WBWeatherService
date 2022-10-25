package main

import (
	"WeatherService/database"
	"WeatherService/models/owRespModel"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	logFile              = "logs.txt"                           //file for logging
	citiesList           = "..\\BashScripts\\Cities.txt"        //List of cities, which we are interested in
	cityDataLoader       = "..\\BashScripts\\GetCityData.sh"    //Bash script generates file with base city information
	cityLoaderResults    = "..\\RequestResults\\CityData\\"     //Generated files with cities can be found here
	weatherDataLoader    = "..\\BashScripts\\GetWeatherData.sh" //Bash script generates file with weather in city information
	weatherLoaderResults = "..\\RequestResults\\WeatherData\\"  //Generated files with weather can be found here
	apikey               = "4ac6651b8e1eef832d0826af46d01b84"
)

/*
------------------------------------------------------------
---------------------	Cities	----------------------------
------------------------------------------------------------
*/
func WriteCitiesToDB(errCh chan error) {
	err, cities := LoadCitiesList()
	if err != nil {
		log.Errorf("Error while reading cities list")
		return
	}

	wg := sync.WaitGroup{}
	defer wg.Wait()
	for _, city := range cities {
		tstr := strings.Split(city, " ")
		go GWriteCity(tstr[0], tstr[1], &wg, errCh)
		wg.Add(1)
	}
}

func LoadCitiesList() (error, []string) {
	data, err := ioutil.ReadFile(citiesList)
	if err != nil {
		log.Errorf("Can't read cities list")
		return err, nil
	}
	strs := strings.Split(string(data), "\r\n")
	return err, strs
}

func GWriteCity(name, country string, wg *sync.WaitGroup, errCh chan error) {
	defer wg.Done()

	/*cmd := exec.Command("powershell",
		cityDataLoader, name, country)

	err := cmd.Run()
	if err != nil {
		log.Errorf("Can't start city shell script execution")
		errCh <- err
		return
	}

	for {
		file, err := os.Stat(cityLoaderResults + name + "_" + country + ".txt")
		if err == nil && time.Since(file.ModTime()) < 30*time.Second {
			break
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}

	var data []byte
	data, err = ioutil.ReadFile(cityLoaderResults + name + "_" + country + ".txt")
	if err != nil {
		log.Errorf("Can't read city response file")
		errCh <- err
		return
	}*/
	url := "http://api.openweathermap.org/geo/1.0/direct?q=" + name + "," +
		country + "&limit=1&appid=" + apikey
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Can't get response from foreign service")
		errCh <- err
		return
	}

	var city []owRespModel.CityResponse
	err = json.NewDecoder(resp.Body).Decode(&city)
	if err != nil {
		log.Errorf("Can't unmarshal city response")
		errCh <- err
		return
	}
	log.Info("Pushing  ", name, " ", country, " to DB")
	database.PostCity(&city[0], errCh)
	return
}

/*
------------------------------------------------------------
---------------------- Weather -----------------------------
------------------------------------------------------------
*/

func WriteWeatherToDB(errCh chan error) {
	err, cities := database.GetCitiesList()
	if err != nil {
		log.Errorf("Can't get cities list from DB")
		errCh <- err
		return
	}

	wg := sync.WaitGroup{}
	defer wg.Wait()

	for _, city := range cities {
		go GWriteWeather(city, &wg, errCh)
		wg.Add(1)
	}
}

func AsyncWeatherUpdate(period time.Duration, closeCh <-chan struct{}, errCh chan error) {
	startTime := time.Now()
	for {
		select {
		case _, ok := <-closeCh:
			if !ok {
				log.Info("Shutting down async updater")
				return
			}
		default:
			if time.Since(startTime) >= period {
				startTime = time.Now()
				log.Info("Async updating weather")
				WriteWeatherToDB(errCh)
			}
		}
	}
}

func GWriteWeather(city owRespModel.CityResponse, wg *sync.WaitGroup, errCh chan error) {
	defer wg.Done()
	/*cmd := exec.Command("powershell",
		weatherDataLoader,
		strconv.FormatFloat(city.Lat, 'f', -1, 64),
		strconv.FormatFloat(city.Lon, 'f', -1, 64),
		city.Name, city.Country)

	err := cmd.Run()
	if err != nil {
		log.Errorf("Can't start weather shell script execution")
		errCh <- err
		return
	}

	for {
		file, err := os.Stat(weatherLoaderResults + city.Name + "_" + city.Country + ".txt")
		if err == nil && time.Since(file.ModTime()) < 30*time.Second {
			break
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
	*/
	log.Info("Getting weather of ", city.Name, " ", city.Country, " from foreign service")
	url := "http://api.openweathermap.org/data/2.5/forecast?" + "lat=" +
		strconv.FormatFloat(city.Lat, 'f', -1, 64) + "&lon=" +
		strconv.FormatFloat(city.Lon, 'f', -1, 64) + "&appid=" + apikey
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Can't get response from foreign service")
		errCh <- err
		return
	}

	var data []byte
	data, err = io.ReadAll(resp.Body)
	f, err := os.OpenFile(weatherLoaderResults+city.Name+"_"+city.Country+".txt",
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Errorf("Can't create or open response weather file")
		errCh <- err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		log.Errorf("Can't write to weather file")
		errCh <- err
	}

	var weather owRespModel.WeatherResponse
	err = json.Unmarshal(data, &weather)
	if err != nil {
		log.Errorf("Can't unmarshal weather response")
		errCh <- err
		return
	}
	weather.City.Id = city.Id
	log.Info("Pushing weather of ", city.Name, " ", city.Country, " to DB")
	database.UpdateWeather(&weather, errCh)
}

func ReadWeatherFormFile(city, country string, date int) (interface{}, error) {
	//data, err := ioutil.ReadFile(weatherLoaderResults + city + "_" + country + ".txt")
	/*if err != nil {
		log.Errorf("Can't read weather response")
		return nil, err
	}*/
	file, err := os.OpenFile(weatherLoaderResults+city+"_"+country+".txt", os.O_RDONLY, 0666)

	var weather owRespModel.WeatherResponse
	err = json.NewDecoder(file).Decode(&weather)
	//err = json.Unmarshal(data, &weather)
	if err != nil {
		log.Errorf("Can't unmarshal weather response")
		return nil, err
	}

	var suitableDateInd int
	diff := absInt(date - weather.List[0].Dt)
	for i, w := range weather.List {
		t := absInt(date - w.Dt)
		if diff > t {
			suitableDateInd = i
			diff = t
		}
	}

	weather.List = weather.List[suitableDateInd : suitableDateInd+1]

	return &weather, err
}

func absInt(a int) int {
	if a < 0 {
		a *= -1
	}
	return a
}
