package infrastructure

import (
	"github.com/diogocardoso/go/lab_1/controllers"
	"github.com/gorilla/mux"
)

func NewRouter(weatherController *controllers.WeatherController) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/weather", weatherController.GetWeather).Methods("GET")
	return router
}
