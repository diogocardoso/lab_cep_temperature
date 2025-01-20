package services

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var WeatherAPIKey = "97b080f3a71949b39d9155111251901"

func TestGetLocationByCEP_ValidCEP(t *testing.T) {
	cep := "01251080"
	expectedLocation := "São Paulo,SP"

	// Cria um servidor de teste para simular a API de localização
	locationServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ws/"+cep+"/json/" {
			response := `{"localidade": "São Paulo", "uf": "SP", "erro": false}`
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
		} else {
			http.NotFound(w, r)
		}
	}))

	defer locationServer.Close()

	service := NewWeatherService(WeatherAPIKey)

	// Executa o teste
	location, cityName, err := service.GetLocationByCEP(cep)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if location != expectedLocation {
		t.Errorf("Expected location '%s', got '%s'", expectedLocation, location)
	}
	if cityName != "São Paulo" {
		t.Errorf("Expected city name 'São Paulo', got '%s'", cityName)
	}
}

func TestGetWeatherByLocation_ValidLocation(t *testing.T) {
	location := "São Paulo,SP"
	expectedTempC := 25.0

	// Cria um servidor de teste para simular a API do clima
	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{"current": {"temp_c": 25.0}}`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer weatherServer.Close()

	// Cria uma instância do WeatherService com a chave da API fictícia
	service := NewWeatherService(WeatherAPIKey)

	// Executa o teste
	weather, err := service.GetWeatherByLocation(location)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if weather == nil {
		t.Fatal("Expected weather to be non-nil")
	}
	if weather.TempC != expectedTempC {
		t.Errorf("Expected TempC to be %v, got %v", expectedTempC, weather.TempC)
	}
}

func TestGetWeatherByLocation_InvalidResponse(t *testing.T) {
	location := "São Paulo,SP"

	// Cria um servidor de teste para simular uma resposta inválida da API do clima
	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"invalid": "response"}`)) // Resposta inválida
	}))
	defer weatherServer.Close()

	// Cria uma instância do WeatherService com a chave da API fictícia
	service := NewWeatherService(WeatherAPIKey)

	// Executa o teste
	_, err := service.GetWeatherByLocation(location)
	if err == nil {
		t.Fatal("Expected error for invalid response, got none")
	}
}
