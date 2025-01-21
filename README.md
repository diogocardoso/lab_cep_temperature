# lab_cep_temperature

Este projeto é um serviço web baseado em Go que fornece informações de temperatura para um determinado CEP. utilizando o [ViaCEP API](https://viacep.com.br/) para obter informações de localização e a [WeatherAPI](https://www.weatherapi.com/) para buscar dados de temperatura.

Deploy : https://cep-temperature-20655388484.us-central1.run.app/weather?cep=72587-260

## Features

- Obter informações de localização usando CEP (Código Postal Brasileiro)
- Recuperar a temperatura atual para a localização
- Converter temperatura para Celsius, Fahrenheit e Kelvin
- Configurável usando variáveis ​​de ambiente

## Prerequisites

- Go 1.21.3 or later
- WeatherAPI API key (sign up at https://www.weatherapi.com)

## Install

1. Clone the repository:
   ```
   git clone https://github.com/renanmav/GoExpert-CEPTemperature-GCR.git
   cd GoExpert-CEPTemperature-GCR
   ```

2. Configure port and api key:
    ```
    config.yaml
    ```

3. Start:

   docker compose up 

