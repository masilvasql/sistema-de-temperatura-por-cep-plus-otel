package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/masilvasql/sistema-de-temperatura-por-cep/configs"
	"github.com/masilvasql/sistema-de-temperatura-por-cep/pkg"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type WeatherUsecase interface {
	GetWeatherByCep(input WehaterInput) (WeaherOutput, error)
}

var ErrorInvalizZipCode = fmt.Errorf("Invalid Zip Code")
var ErrorZipCodeNotFound = fmt.Errorf("can not find zipcode")

type weatherUsecase struct {
	EnvConfig *configs.Config
}

type WehaterInput struct {
	Cep string `json:"cep"`
}

type WeaherOutput struct {
	TemperatureInCelsius    float64 `json:"temp_C"`
	TemperatureInFahrenheit float64 `json:"temp_F"`
	TemperatureInKelvin     float64 `json:"temp_K"`
	City                    string  `json:"city"`
}

type WeatherAPIResponse struct {
	Current Current `json:"current"`
	City    string  `json:"city"`
}

type Current struct {
	TempC float64 `json:"temp_c"`
}

func NewWeatherUsecase(envConfig *configs.Config) WeatherUsecase {
	return &weatherUsecase{
		EnvConfig: envConfig,
	}
}

func (w *weatherUsecase) GetWeatherByCep(input WehaterInput) (WeaherOutput, error) {

	if !pkg.IsValidZipCode(input.Cep) {
		return WeaherOutput{}, ErrorInvalizZipCode
	}

	weatherAPIResponse, err := w.doCepRequest(input)
	if err != nil {
		return WeaherOutput{}, err
	}

	weatherOutput := WeaherOutput{
		TemperatureInCelsius:    weatherAPIResponse.Current.TempC,
		TemperatureInFahrenheit: (weatherAPIResponse.Current.TempC * 9 / 5) + 32,
		TemperatureInKelvin:     weatherAPIResponse.Current.TempC + 273.15,
		City:                    weatherAPIResponse.City,
	}

	return weatherOutput, nil
}

func (w *weatherUsecase) doCepRequest(input WehaterInput) (WeatherAPIResponse, error) {

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	jsonBody, err := json.Marshal(input)
	if err != nil {
		return WeatherAPIResponse{}, err
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", "http://go_service_b:8081/weather", bytes.NewBuffer(jsonBody))

	if err != nil {
		return WeatherAPIResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return WeatherAPIResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return WeatherAPIResponse{}, ErrorZipCodeNotFound
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherAPIResponse{}, err
	}

	var response WeatherAPIResponse
	err = json.Unmarshal(body, &response)

	if err != nil {
		return WeatherAPIResponse{}, err
	}

	return response, nil

}
