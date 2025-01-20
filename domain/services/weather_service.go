package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/diogocardoso/go/lab_1/domain/models"
)

type WeatherService interface {
	GetLocationByCEP(cep string) (string, string, error)
	GetWeatherByCEP(cep string) (*models.Weather, error)
	GetWeatherByLocation(location string) (*models.Weather, error)
}

type weatherService struct {
	weatherAPIKey string
}

func NewWeatherService(apiKey string) WeatherService {
	return &weatherService{weatherAPIKey: apiKey}
}

func (s *weatherService) GetWeatherByCEP(cep string) (*models.Weather, error) {
	log.Printf("Iniciando a obtenção do clima para o CEP: %s", cep)

	if len(cep) != 8 || !isNumeric(cep) {
		err := errors.New("CEP inválido")
		log.Printf("Erro: %s", err)
		return nil, err
	}

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
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Erro ao fazer requisição para o CEP %s: %s", cep, err)
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := errors.New("não foi possível encontrar o CEP")
		log.Printf("Erro: %s", err)
		return "", "", err
	}

	var result struct {
		Localidade string `json:"localidade"`
		Uf         string `json:"uf"`
		Erro       bool   `json:"erro"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Erro ao decodificar resposta para o CEP %s: %s", cep, err)
		return "", "", err
	}

	if result.Erro {
		err := errors.New("não foi possível encontrar o CEP")
		log.Printf("Erro: %s", err)
		return "", "", err
	}

	location := fmt.Sprintf("%s,%s", result.Localidade, result.Uf)
	cityName := result.Localidade
	log.Printf("Localização obtida para o CEP %s: %s", cep, location)
	return location, cityName, nil
}

func (s *weatherService) GetWeatherByLocation(location string) (*models.Weather, error) {
	log.Printf("Obtendo clima para a localização: %s", location)
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", s.weatherAPIKey, location)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Erro ao fazer requisição para a localização %s: %s", location, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := errors.New("falha ao recuperar dados do clima")
		log.Printf("Erro: %s", err)
		log.Printf("Erro URL: %s", url)
		return nil, err
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
