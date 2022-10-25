package main

import (
	"WeatherService/database"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetListCities(c *gin.Context) {
	log.Info("Request: al cities")
	err, cities := database.GetCitiesList()
	if err != nil {
		notFound(c, err, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, cities)
}

func GetCityShortPrediction(c *gin.Context) {
	log.Info("Short prediction request: ", c.Query("city"), c.Query("country"))
	err, resp := database.GetCityWeather(c.Query("city"), c.Query("country"))
	if err != nil {
		notFound(c, err, "City not found")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetCityLongPrediction(c *gin.Context) {
	log.Info("Long prediction request: ", c.Query("city"), c.Query("country"), c.Query("date"))
	date, err := strconv.Atoi(c.Query("date"))
	if err != nil {
		badRequest(c, err, "Can't parse user's date")
		return
	}
	resp, err := ReadWeatherFormFile(c.Query("city"), c.Query("country"), date)
	if err != nil {
		notFound(c, err, "City not found")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func notFound(c *gin.Context, err error, desc string) {
	c.String(http.StatusNotFound, desc)
	log.Errorf("%s. %s", desc, err)
}

func badRequest(c *gin.Context, err error, desc string) {
	c.String(http.StatusBadRequest, desc)
	log.Errorf("%s. %s", desc, err)
}

func ServiceRun() {
	router := gin.Default()

	router.GET("/fullforecast/", GetCityLongPrediction)
	router.GET("/forecast/", GetCityShortPrediction)
	router.GET("/cities", GetListCities)

	err := router.Run(":8080")
	if err != nil {
		log.Errorf("Can't start server: %s", err)
		return
	}
	log.Info("Server running")
}
