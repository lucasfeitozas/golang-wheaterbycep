package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherByCEPHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		method         string
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "CEP válido - São Paulo",
			path:           "/weatherbycep/01310100",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedMsg:    "",
		},
		{
			name:           "CEP válido com hífen",
			path:           "/weatherbycep/01310-100",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedMsg:    "",
		},
		{
			name:           "CEP inválido - formato incorreto",
			path:           "/weatherbycep/123",
			method:         "GET",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedMsg:    "invalid zipcode",
		},
		{
			name:           "CEP não encontrado",
			path:           "/weatherbycep/00000000",
			method:         "GET",
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "can not find zipcode",
		},
		{
			name:           "CEP não fornecido",
			path:           "/weatherbycep/",
			method:         "GET",
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "cep parameter is required",
		},
		{
			name:           "Método não permitido - POST",
			path:           "/weatherbycep/01310100",
			method:         "POST",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedMsg:    "method not allowed",
		},
		{
			name:           "Endpoint não encontrado",
			path:           "/invalid",
			method:         "GET",
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "endpoint not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(weatherByCEPHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler retornou status code errado: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				// Verifica se a resposta é um JSON válido com dados de temperatura
				var weather WeatherData
				if err := json.Unmarshal(rr.Body.Bytes(), &weather); err != nil {
					t.Errorf("Resposta não é um JSON válido: %v", err)
				}

				// Verifica se os campos de temperatura estão presentes
				if weather.TempC == 0 && weather.TempF == 0 && weather.TempK == 0 {
					t.Errorf("Dados de temperatura não encontrados na resposta")
				}

				// Verifica conversões de temperatura
				expectedTempF := (weather.TempC * 9 / 5) + 32
				expectedTempK := weather.TempC + 273.15

				if weather.TempF != expectedTempF {
					t.Errorf("Conversão Fahrenheit incorreta: got %v want %v",
						weather.TempF, expectedTempF)
				}

				if weather.TempK != expectedTempK {
					t.Errorf("Conversão Kelvin incorreta: got %v want %v",
						weather.TempK, expectedTempK)
				}
			} else if tt.expectedMsg != "" {
				// Verifica mensagem de erro
				var errorResp ErrorResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &errorResp); err != nil {
					t.Errorf("Resposta de erro não é um JSON válido: %v", err)
				}

				if errorResp.Message != tt.expectedMsg {
					t.Errorf("Mensagem de erro incorreta: got %v want %v",
						errorResp.Message, tt.expectedMsg)
				}
			}
		})
	}
}

func TestIsValidCEP(t *testing.T) {
	tests := []struct {
		cep      string
		expected bool
	}{
		{"01310100", true},
		{"01310-100", true},
		{"12345678", true},
		{"123", false},
		{"1234567890", false},
		{"abcd1234", false},
		{"", false},
		{"123-456", false},
		{"12.345.678", false},
	}

	for _, tt := range tests {
		t.Run(tt.cep, func(t *testing.T) {
			result := isValidCEP(tt.cep)
			if result != tt.expected {
				t.Errorf("isValidCEP(%s) = %v, want %v", tt.cep, result, tt.expected)
			}
		})
	}
}

func TestFormatCEP(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"01310-100", "01310100"},
		{"01310 100", "01310100"},
		{"01310100", "01310100"},
		{"123-45-678", "12345678"},
		{"12 34 56 78", "12345678"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := formatCEP(tt.input)
			if result != tt.expected {
				t.Errorf("formatCEP(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark para testar performance
func BenchmarkWeatherByCEPHandler(b *testing.B) {
	req, _ := http.NewRequest("GET", "/weatherbycep/01310100", nil)

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(weatherByCEPHandler)
		handler.ServeHTTP(rr, req)
	}
}
