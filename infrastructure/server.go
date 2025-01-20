package infrastructure

import (
	"log"
	"net/http"

	"github.com/diogocardoso/go/lab_1/configs"
	"github.com/diogocardoso/go/lab_1/controllers"
	"github.com/diogocardoso/go/lab_1/domain/services"
	"github.com/gorilla/mux"
)

type Server struct {
	Config            *configs.Config
	WeatherService    services.WeatherService
	WeatherController *controllers.WeatherController
	Router            *mux.Router
}

func NewServer(config *configs.Config) *Server {
	weatherService := services.NewWeatherService(config.WeatherAPIKey)
	weatherController := controllers.NewWeatherController(weatherService)
	router := NewRouter(weatherController)

	return &Server{
		Config:            config,
		WeatherService:    weatherService,
		WeatherController: weatherController,
		Router:            router,
	}
}

func (s *Server) Start() {
	log.Printf("Starting server on port %s...", s.Config.Port)
	if err := http.ListenAndServe(":"+s.Config.Port, s.Router); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
