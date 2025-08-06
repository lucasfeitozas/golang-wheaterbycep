package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// httpClient é um cliente HTTP personalizado com configuração TLS tolerante para Cloud Run
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false, // Mantém a verificação de certificado
			MinVersion:         tls.VersionTLS12,
		},
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
		ForceAttemptHTTP2:  true,
	},
}

// CEPData representa a estrutura de dados retornada pela API do ViaCEP
type CEPData struct {
	CEP         string      `json:"cep"`
	Logradouro  string      `json:"logradouro"`
	Complemento string      `json:"complemento"`
	Bairro      string      `json:"bairro"`
	Localidade  string      `json:"localidade"`
	UF          string      `json:"uf"`
	IBGE        string      `json:"ibge"`
	GIA         string      `json:"gia"`
	DDD         string      `json:"ddd"`
	SIAFI       string      `json:"siafi"`
	Erro        interface{} `json:"erro,omitempty"`
}

// WeatherData representa a estrutura de dados de temperatura
type WeatherData struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// ErrorResponse representa a estrutura de resposta de erro
type ErrorResponse struct {
	Message string `json:"message"`
}

// isValidCEP valida se o CEP está no formato correto
func isValidCEP(cep string) bool {
	// Remove traços e espaços
	cep = strings.ReplaceAll(cep, "-", "")
	cep = strings.ReplaceAll(cep, " ", "")

	// Verifica se tem 8 dígitos
	if len(cep) != 8 {
		return false
	}

	// Verifica se contém apenas números
	matched, _ := regexp.MatchString(`^\d{8}$`, cep)
	return matched
}

// formatCEP formata o CEP removendo caracteres especiais
func formatCEP(cep string) string {
	cep = strings.ReplaceAll(cep, "-", "")
	cep = strings.ReplaceAll(cep, " ", "")
	return cep
}

// CustomError representa erros customizados com códigos HTTP
type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

// searchCEP faz a consulta na API do ViaCEP
func searchCEP(cep string) (*CEPData, *CustomError) {
	// Valida o CEP
	if !isValidCEP(cep) {
		return nil, &CustomError{Code: 422, Message: "invalid zipcode"}
	}

	// Formata o CEP
	formattedCEP := formatCEP(cep)

	// Monta a URL da API
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", formattedCEP)

	// Faz a requisição HTTP usando o cliente personalizado
	resp, err := httpClient.Get(url)
	if err != nil {
		// Se falhar com HTTPS, tenta com HTTP como fallback
		log.Printf("Erro com HTTPS, tentando HTTP: %v\n", err)
		httpURL := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", formattedCEP)
		resp, err = httpClient.Get(httpURL)
		if err != nil {
			log.Printf("Erro ao fazer requisição para ViaCEP: %v\n", err)
			return nil, &CustomError{Code: 500, Message: "internal server error"}
		}
	}
	defer resp.Body.Close()

	// Verifica se a resposta foi bem-sucedida
	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro na resposta do ViaCEP: %s\n", resp.Status)
		return nil, &CustomError{Code: 500, Message: "internal server error"}
	}

	// Lê o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler o corpo da resposta: %v\n", err)
		return nil, &CustomError{Code: 500, Message: "internal server error"}
	}

	// Decodifica o JSON
	var cepData CEPData
	if err := json.Unmarshal(body, &cepData); err != nil {
		fmt.Printf("Erro ao decodificar JSON: %v\n", err)
		return nil, &CustomError{Code: 500, Message: "internal server error"}
	}

	// Verifica se o CEP foi encontrado
	if cepData.Erro != nil {
		fmt.Printf("CEP não encontrado: %s\n", cep)
		return nil, &CustomError{Code: 404, Message: "can not find zipcode"}
	}

	return &cepData, nil
}

// getWeatherData busca os dados de temperatura usando uma API gratuita
func getWeatherData(city, state string) (*WeatherData, *CustomError) {
	// Forma alternativa: usar wttr.in que é gratuito e não requer chave
	cityFormatted := strings.ReplaceAll(city, " ", "+")
	stateFormatted := strings.ReplaceAll(state, " ", "+")
	location := fmt.Sprintf("%s,%s,Brazil", cityFormatted, stateFormatted)

	// URL da API wttr.in em formato JSON
	url := fmt.Sprintf("https://wttr.in/%s?format=j1", url.QueryEscape(location))

	resp, err := httpClient.Get(url)
	if err != nil {
		fmt.Printf("Erro ao fazer requisição para wttr.in: %v\n", err)
		return nil, &CustomError{Code: 500, Message: "internal server error"}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erro na resposta da API wttr.in: %s\n", resp.Status)
		return nil, &CustomError{Code: 500, Message: "internal server error"}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler o corpo da resposta: %v\n", err)
		return nil, &CustomError{Code: 500, Message: "internal server error"}
	}

	// Estrutura específica para wttr.in
	var wttrResponse struct {
		CurrentCondition []struct {
			TempC string `json:"temp_C"`
		} `json:"current_condition"`
	}

	if err := json.Unmarshal(body, &wttrResponse); err != nil {
		fmt.Printf("Erro ao decodificar JSON: %v\n", err)
		return nil, &CustomError{Code: 500, Message: "internal server error"}
	}

	if len(wttrResponse.CurrentCondition) == 0 {
		fmt.Println("Dados climáticos não disponíveis para a localização fornecida.")
		return nil, &CustomError{Code: 500, Message: "weather data not available"}
	}

	// Converte temperatura de string para float64
	tempCStr := wttrResponse.CurrentCondition[0].TempC
	tempC, err := strconv.ParseFloat(tempCStr, 64)
	if err != nil {
		fmt.Printf("Erro ao converter temperatura: %v\n", err)
		return nil, &CustomError{Code: 500, Message: "internal server error"}
	}

	// Calcula as conversões de temperatura
	tempF := (tempC * 9 / 5) + 32 // Celsius para Fahrenheit
	tempK := tempC + 273.15       // Celsius para Kelvin

	return &WeatherData{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}

// weatherByCEPHandler lida com as requisições GET para /weatherbycep/{cep}
func weatherByCEPHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica se é um GET
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "method not allowed"})
		return
	}

	// Extrai o CEP do path da URL
	// Remove o prefixo "/weatherbycep/" para obter o CEP
	path := r.URL.Path
	if !strings.HasPrefix(path, "/weatherbycep/") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "endpoint not found"})
		return
	}

	cep := strings.TrimPrefix(path, "/weatherbycep/")
	if cep == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "cep parameter is required"})
		return
	}

	// Busca os dados do CEP
	cepData, cepErr := searchCEP(cep)
	if cepErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(cepErr.Code)
		json.NewEncoder(w).Encode(ErrorResponse{Message: cepErr.Message})
		return
	}

	// Busca dados climáticos
	weather, weatherErr := getWeatherData(cepData.Localidade, cepData.UF)
	if weatherErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(weatherErr.Code)
		json.NewEncoder(w).Encode(ErrorResponse{Message: weatherErr.Message})
		return
	}

	// Retorna os dados de temperatura em caso de sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weather)
}

func main() {
	// Configura o handler para o endpoint /weatherbycep/{cep}
	http.HandleFunc("/weatherbycep/", weatherByCEPHandler)

	// Define a porta do servidor
	port := ":8080"

	fmt.Printf("🌡️  Servidor iniciado na porta %s\n", port)
	fmt.Println("📡 Endpoint disponível: GET /weatherbycep/{cep}")
	fmt.Println("📋 Exemplo de uso: GET /weatherbycep/01310100")

	// Inicia o servidor
	log.Fatal(http.ListenAndServe(port, nil))
}
