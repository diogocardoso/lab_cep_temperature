package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/diogocardoso/go/lab_1/domain/models"
)

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

type WeatherService interface {
	GetWeatherByCEP(cep string) (*models.Weather, error)
	GetWeatherByLocation(location string) (*models.Weather, error)
	GetLocationByCEP(cep string) (string, string, error)
}

type weatherService struct {
	weatherAPIKey string
}

func NewWeatherService(apiKey string) WeatherService {
	return &weatherService{weatherAPIKey: apiKey}
}

func (s *weatherService) GetWeatherByCEP(cep string) (*models.Weather, error) {
	location, cityName, err := s.GetLocationByCEP(cep)
	if err != nil {
		log.Printf("Erro ao obter localização para o CEP %s: %s", cep, err)
		return nil, err
	}

	weather, err := s.GetWeatherByLocation(location)
	if err != nil {
		log.Printf("Erro ao obter clima para a localização %s: %s", location, err)
		return nil, err
	}

	weather.CityName = cityName
	log.Printf("Clima obtido com sucesso para o CEP %s: %+v", cep, weather)
	return weather, nil
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (s *weatherService) GetLocationByCEP(cep string) (string, string, error) {
	log.Printf("Obtendo localização para o CEP: %s", cep)

	if len(cep) != 8 && len(cep) != 9 {
		err := errors.New("invalid zipcode")
		log.Printf("Erro: %s, CEP: %s", err, cep)
		return "", "", &HTTPError{StatusCode: 422, Message: err.Error()}
	}

	if strings.Contains(cep, "-") {
		cep = strings.ReplaceAll(cep, "-", "")
		log.Printf("CEP: %s", cep)
	}

	if !isNumeric(cep) {
		err := errors.New("invalid zipcode")
		log.Printf("Erro: %s", err)
		return "", "", &HTTPError{StatusCode: 422, Message: err.Error()}
	}

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	log.Printf("Obtendo dados em: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Erro ao fazer requisição para o CEP %s: %s", cep, err)
		return "", "", &HTTPError{StatusCode: 404, Message: "can not find zipcode"}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := errors.New("não foi possível encontrar o CEP")
		log.Printf("Erro: %s", err)
		return "", "", &HTTPError{StatusCode: 404, Message: "can not find zipcode"}
	}

	var result struct {
		Localidade string `json:"localidade"`
		Uf         string `json:"uf"`
		Erro       bool   `json:"erro"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Erro ao decodificar resposta para o CEP %s: %s", cep, err)
		return "", "", &HTTPError{StatusCode: 404, Message: "can not find zipcode"}
	}

	if result.Erro {
		err := errors.New("não foi possível encontrar o CEP")
		log.Printf("Erro: %s", err)
		return "", "", &HTTPError{StatusCode: 404, Message: "can not find zipcode"}
	}

	location := fmt.Sprintf("%s,%s", result.Localidade, result.Uf)
	cityName := result.Localidade
	log.Printf("Localização obtida para o CEP %s: %s", cep, location)
	return location, cityName, nil
}

func (s *weatherService) GetWeatherByLocation(location string) (*models.Weather, error) {
	encodedCity := url.QueryEscape(location)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", s.weatherAPIKey, encodedCity)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro URL: %s", url)
		return nil, &HTTPError{StatusCode: 404, Message: fmt.Sprintf("error getting city (%s) weather data", location)}
	}

	var result struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Erro ao decodificar resposta para a localização %s: %s", location, err)
		return nil, err
	}

	tempC := result.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15

	weather := &models.Weather{TempC: tempC, TempF: tempF, TempK: tempK}
	log.Printf("Clima obtido para a localização %s: %+v", location, weather)
	return weather, nil
}
