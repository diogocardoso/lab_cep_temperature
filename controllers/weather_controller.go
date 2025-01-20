package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/diogocardoso/go/lab_1/domain/services"
)

type WeatherController struct {
	WeatherService services.WeatherService
}

func NewWeatherController(ws services.WeatherService) *WeatherController {
	return &WeatherController{WeatherService: ws}
}

func (wc *WeatherController) GetWeather(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if len(cep) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	weather, err := wc.WeatherService.GetWeatherByCEP(cep)
	if err != nil {
		if err.Error() == "invalid zipcode" {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		} else if err.Error() == "can not find zipcode" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	response, err := json.Marshal(weather)
	if err != nil {
		http.Error(w, "error processing response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
