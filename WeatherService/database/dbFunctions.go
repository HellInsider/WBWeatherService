package database

import (
	"WeatherService/models/owRespModel"
	sqlFunctions "WeatherService/sqlFunctions"
	"database/sql"
	"fmt"
	"github.com/bradfitz/slice"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const sqlConnect string = "user=postgres password=admin dbname=WB_weather sslmode=disable"

func PostCity(city *owRespModel.CityResponse, errCh chan error) {
	db, err := sql.Open("postgres", sqlConnect)
	if err != nil {
		log.Errorf("Error while opening DB: %s", err)
		errCh <- err
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		errCh <- err
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(sqlFunctions.PostCity, city.Name, city.Country, city.Lat, city.Lon)
	if err != nil {
		log.Errorf("Error while execution sqlFunctions.PostCity(): %s", err)
		errCh <- err
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Errorf("Error while commit to Cities: %s", err)
		errCh <- err
		return
	}

	return
}

func UpdateWeather(weather *owRespModel.WeatherResponse, errCh chan error) {
	db, err := sql.Open("postgres", sqlConnect)
	if err != nil {
		log.Errorf("Error while opening DB: %s", err)
		errCh <- err
		return
	}
	defer db.Close()

	for _, w := range weather.List {
		tx, err := db.Begin()
		//_, err = tx.Exec("select postweather($1,$2,$3,$4)", weather.City.Id, w.Dt, w.Main.Temp, w.DtTxt)
		_, err = tx.Exec(sqlFunctions.PostWeather, weather.City.Id, w.Dt, w.Main.Temp, w.DtTxt)
		if err != nil {
			log.Errorf("Error while execution sqlFunctions.PostWeather(): %s", err)
			errCh <- err
			return
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			log.Errorf("Error while commit to Weather: %s", err)
			errCh <- err
			return
		}
	}

	return
}

func GetCitiesList() (error, []owRespModel.CityResponse) {
	db, err := sql.Open("postgres", sqlConnect)
	if err != nil {
		log.Errorf("Connection to DB error: %s", err)
		return err, nil
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Error(fmt.Errorf("transaction error: %s", err))
		return err, nil
	}
	defer tx.Rollback()

	orderRows, err := db.Query(`select city_id, city_name, country, lat, lon from public."Cities"`)
	if err != nil {
		log.Errorf("error getting rows from Cities: %s", err)
		return err, nil
	}
	defer orderRows.Close()

	var data owRespModel.CityResponse
	var res []owRespModel.CityResponse
	for orderRows.Next() {
		err = orderRows.Scan(&data.Id, &data.Name, &data.Country, &data.Lat, &data.Lon)
		if err != nil {
			log.Errorf("error parsing rows from Cities: %s", err)
			return err, nil
		}
		res = append(res, data)
	}
	return err, res
}

func GetCityWeather(name, country string) (error, interface{}) {
	type weather struct {
		Dt    int    `json:"dt"`
		DtTxt string `json:"dt_txt"`
	}

	type resp struct {
		Name       string  `json:"name"`
		Country    string  `json:"country"`
		AvTemp     float64 `json:"av_temp"`
		WeatherArr []weather
	}

	db, err := sql.Open("postgres", sqlConnect)
	if err != nil {
		log.Errorf("Connection to DB error: %s", err)
		return err, nil
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Error(fmt.Errorf("transaction error: %s", err))
		return err, nil
	}
	defer tx.Rollback()

	var rResp resp
	rResp.Name = name
	rResp.Country = country

	orderRows, err := db.Query(sqlFunctions.GetCityWeather, name, country)
	if err != nil {
		log.Errorf("error getting rows from Weather: %s", err)
		return err, nil
	}
	defer orderRows.Close()

	var (
		w           weather
		i           int
		temperature float64
		_temp       float64
	)
	for orderRows.Next() {
		err = orderRows.Scan(&_temp, &w.Dt, &w.DtTxt)
		if err != nil {
			log.Errorf("error parsing rows from Weather: %s", err)
			return err, nil
		}
		rResp.WeatherArr = append(rResp.WeatherArr, w)
		temperature += _temp
		i++
	}

	slice.Sort(rResp.WeatherArr[:], func(i, j int) bool {
		return rResp.WeatherArr[i].Dt < rResp.WeatherArr[j].Dt
	})
	rResp.AvTemp = temperature / float64(i)
	return err, rResp
}
