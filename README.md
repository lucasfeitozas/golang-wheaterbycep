# API de Consulta de CEP e Temperatura

Uma API REST em Go para consultar dados de CEP usando a API do ViaCEP e obter informações de temperatura da localização.

## 🚀 Teste rápido

**Teste a API diretamente no Cloud Run:**
```bash
curl -X GET https://golang-weatherbycep-390503355828.us-central1.run.app/weatherbycep/69086129
```

## 🚀 Funcionalidades

- ✅ API REST com endpoint `GET /weatherbycep/{cep}`
- ✅ Consulta de CEP via API do ViaCEP
- ✅ Busca automática de temperatura do local
- ✅ Validação de formato de CEP
- ✅ Suporte a CEP com ou sem hífen
- ✅ Conversões de temperatura (Celsius, Fahrenheit, Kelvin)
- ✅ Tratamento de erros com códigos HTTP apropriados
- ✅ Resposta em formato JSON

## 📋 Pré-requisitos

- Go 1.19 ou superior
- Conexão com a internet

## 🛠️ Instalação e Execução

1. Clone ou baixe este repositório
2. Navegue até o diretório do projeto:
   ```bash
   cd golang-weatherbycep
   ```
3. Execute o servidor:
   ```bash
   go run main.go
   ```
4. O servidor estará disponível em `http://localhost:8080`

## 🐳 Execução com Docker

### Usando Docker Compose (recomendado):
```bash
# Build e execução
docker-compose up --build

# Execução em background
docker-compose up -d

# Parar os serviços
docker-compose down
```

### Usando Docker diretamente:
```bash
# Build da imagem
docker build -t weatherbycep .

# Execução do container
docker run -p 8080:80 weatherbycep

# Execução em background
docker run -d -p 8080:80 --name weatherbycep weatherbycep
```

A aplicação estará disponível em `http://localhost:8080`

## 🧪 Testes Automatizados

### Executar todos os testes:
```bash
# Testes unitários
go test -v

# Testes com coverage
go test -cover

# Benchmark de performance
go test -bench=.

# Usando o script de teste
chmod +x test.sh
./test.sh
```

### Testes incluídos:
- ✅ Validação de CEP (formato correto/incorreto)
- ✅ Teste de endpoints com CEPs válidos
- ✅ Teste de códigos de erro HTTP
- ✅ Verificação de conversões de temperatura
- ✅ Teste de métodos HTTP não permitidos
- ✅ Benchmark de performance
- ✅ Teste de formatação de CEP

### Exemplo de execução dos testes:
```bash
$ go test -v
=== RUN   TestWeatherByCEPHandler
=== RUN   TestWeatherByCEPHandler/CEP_válido_-_São_Paulo
=== RUN   TestWeatherByCEPHandler/CEP_inválido_-_formato_incorreto
=== RUN   TestWeatherByCEPHandler/CEP_não_encontrado
=== RUN   TestWeatherByCEPHandler/Método_não_permitido_-_POST
--- PASS: TestWeatherByCEPHandler (2.34s)
=== RUN   TestIsValidCEP
--- PASS: TestIsValidCEP (0.00s)
=== RUN   TestFormatCEP
--- PASS: TestFormatCEP (0.00s)
PASS
```

## ⚡ Comandos Úteis (Makefile)

```bash
# Ver todos os comandos disponíveis
make help

# Build da aplicação
make build

# Executar localmente
make run

# Executar todos os testes
make test-all

# Build e execução com Docker
make docker-compose-up

# Parar containers
make docker-compose-down

# Limpar arquivos temporários
make clean
```

## 🚀 CI/CD

Este projeto inclui pipeline de CI/CD com GitHub Actions que:

- ✅ Executa testes automatizados
- ✅ Verifica coverage de código
- ✅ Executa benchmarks de performance
- ✅ Build da aplicação
- ✅ Build e teste da imagem Docker
- ✅ Validação em múltiplas versões do Go

O pipeline é executado automaticamente em:
- Push para branches `main` e `develop`
- Pull requests para `main`

## 🎯 Como usar

### Endpoint disponível:
```
GET /weatherbycep/{cep}
```

### 🌐 Teste direto no Cloud Run:
```bash
curl -X GET https://golang-weatherbycep-390503355828.us-central1.run.app/weatherbycep/69086129
```

### Exemplos de uso com curl:

**CEP de Manaus (teste recomendado):**
```bash
curl -X GET https://golang-weatherbycep-390503355828.us-central1.run.app/weatherbycep/69086129
```

**CEP de São Paulo:**
```bash
curl -X GET https://golang-weatherbycep-390503355828.us-central1.run.app/weatherbycep/01310100
```

**CEP com hífen:**
```bash
curl -X GET https://golang-weatherbycep-390503355828.us-central1.run.app/weatherbycep/01310-100
```

## 📊 Respostas da API

### ✅ Sucesso (200 OK)
```json
{
  "temp_C": 17,
  "temp_F": 62.6,
  "temp_K": 290.15
}
```

### ❌ CEP inválido (422 Unprocessable Entity)
```json
{
  "message": "invalid zipcode"
}
```

### ❌ CEP não encontrado (404 Not Found)
```json
{
  "message": "can not find zipcode"
}
```

### ❌ Método não permitido (405 Method Not Allowed)
```json
{
  "message": "method not allowed"
}
```

### ❌ CEP não fornecido (400 Bad Request)
```json
{
  "message": "cep parameter is required"
}
```

## 🔍 Dados retornados

### Resposta de sucesso:
- **temp_C**: Temperatura em graus Celsius
- **temp_F**: Temperatura em graus Fahrenheit  
- **temp_K**: Temperatura em Kelvin

### Códigos de status HTTP:
- **200**: Sucesso - retorna dados de temperatura
- **400**: CEP não fornecido no path
- **404**: CEP não encontrado
- **405**: Método HTTP não permitido (apenas GET é aceito)
- **422**: CEP com formato inválido
- **500**: Erro interno do servidor

## ⚠️ Tratamento de erros

A API trata os seguintes casos de erro conforme especificado:

- **CEP inválido**: Formato incorreto (não possui 8 dígitos numéricos)
- **CEP não encontrado**: CEP válido mas inexistente na base de dados
- **CEP não fornecido**: Path sem CEP (apenas `/weatherbycep/`)
- **Método não permitido**: Tentativa de usar POST, PUT, DELETE, etc.
- **Problemas de conexão**: Falhas nas APIs externas

## 🌐 APIs utilizadas

Esta aplicação utiliza duas APIs gratuitas:

### ViaCEP (Dados de CEP):
- URL base: `https://viacep.com.br/ws/`
- Formato: `https://viacep.com.br/ws/{CEP}/json/`
- Documentação: [ViaCEP](https://viacep.com.br/)

### wttr.in (Dados Climáticos):
- URL base: `https://wttr.in/`
- Formato: `https://wttr.in/{location}?format=j1`
- Documentação: [wttr.in](https://wttr.in/:help)

## 📝 Estrutura do código

- **main.go**: Arquivo principal com toda a lógica da API
- **main_test.go**: Testes automatizados da aplicação
- **Dockerfile**: Configuração para containerização
- **docker-compose.yml**: Orquestração de containers para desenvolvimento
- **test.sh**: Script para execução de testes
- **Estruturas**: `CEPData`, `WeatherData`, `ErrorResponse`
- **Handlers**: `weatherByCEPHandler` para processar requisições GET
- **Validação**: Função `isValidCEP` para validar formato do CEP
- **APIs externas**: Integração com ViaCEP e wttr.in

## 🤝 Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para:

1. Fazer fork do projeto
2. Criar uma branch para sua feature
3. Commitar suas mudanças
4. Fazer push para a branch
5. Abrir um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.
