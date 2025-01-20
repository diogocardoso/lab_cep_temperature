package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	WeatherAPIKey string
	Port          string
}

func LoadConfig() (*Config, error) {
	// Define o nome e o tipo do arquivo de configuração
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // Diretório atual

	// Tenta ler o arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("erro ao carregar o arquivo de configuração: %w", err)
	}

	// Cria uma instância de Config e preenche com valores do arquivo de configuração
	config := &Config{
		WeatherAPIKey: viper.GetString("WeatherAPIKey"),
		Port:          viper.GetString("Port"),
	}

	// Se a porta não estiver definida, utiliza a porta padrão
	if config.Port == "" {
		config.Port = "8080" // Porta padrão
	}

	return config, nil
}
