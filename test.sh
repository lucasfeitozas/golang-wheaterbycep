#!/bin/bash

echo "🧪 Executando testes da API de CEP e Temperatura"
echo "================================================="

# Executa os testes
echo "📋 Executando testes unitários..."
go test -v

echo ""
echo "📊 Executando testes com coverage..."
go test -cover

echo ""
echo "⚡ Executando benchmark..."
go test -bench=.

echo ""
echo "✅ Testes concluídos!"
