#!/bin/bash

# NOTE: If you get a "Permission denied" error, run:
#   chmod +x scripts/test-pipeline.sh
# Script para testar a pipeline localmente antes de subir para o GitHub
# Usage: ./scripts/test-pipeline.sh

set -e

echo "🔧 Testando Pipeline CI/CD Localmente"
echo "====================================="
echo

# Verificar se estamos no diretório correto
if [ ! -f "go.mod" ]; then
    echo "❌ Erro: Execute este script na raiz do projeto (onde está o go.mod)"
    exit 1
fi

echo "📋 Informações do ambiente:"
echo "  Go version: $(go version)"
echo "  Branch atual: $(git branch --show-current)"
echo "  Último commit: $(git log --oneline -1)"
echo

# 1. Download de dependências
echo "📦 1. Baixando dependências..."
go mod download
echo "✅ Dependências baixadas"
echo

# 2. Testes
echo "🧪 2. Executando testes..."
go test ./... -v -coverprofile=coverage.out
echo "✅ Testes passaram"
echo

# 3. Verificação de formatação
echo "📝 3. Verificando formatação..."
UNFORMATTED=$(gofmt -s -l . | wc -l)
if [ "$UNFORMATTED" -gt 0 ]; then
    echo "❌ Código não está formatado corretamente:"
    gofmt -s -l .
    echo "Para corrigir execute: go fmt ./..."
    exit 1
else
    echo "✅ Código está bem formatado"
fi
echo

# 4. Linter (se estiver instalado)
echo "🔍 4. Executando linter..."
if command -v golangci-lint >/dev/null 2>&1; then
    golangci-lint run --timeout=5m
    echo "✅ Linter passou"
elif [ -f "$HOME/go/bin/golangci-lint" ]; then
    $HOME/go/bin/golangci-lint run --timeout=5m
    echo "✅ Linter passou"
else
    echo "⚠️  golangci-lint não instalado, pulando..."
    echo "Para instalar: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
fi
echo

# 5. Build dos exemplos
echo "🔨 5. Buildando exemplos..."
mkdir -p bin
go build -o bin/basic ./examples/basic/
go build -o bin/capacity-demo ./examples/capacity-demo/
echo "✅ Exemplos buildados com sucesso"
echo

# 6. Teste cross-platform (opcional)
echo "🌍 6. Testando build cross-platform..."
GOOS=linux GOARCH=amd64 go build -o /tmp/test-linux ./examples/basic/ 2>/dev/null && echo "✅ Linux OK" || echo "❌ Linux falhou"
GOOS=windows GOARCH=amd64 go build -o /tmp/test-windows.exe ./examples/basic/ 2>/dev/null && echo "✅ Windows OK" || echo "❌ Windows falhou"
GOOS=darwin GOARCH=amd64 go build -o /tmp/test-macos ./examples/basic/ 2>/dev/null && echo "✅ macOS OK" || echo "❌ macOS falhou"
rm -f /tmp/test-*
echo

# 7. Verificação de cobertura
echo "📊 7. Relatório de cobertura:"
if [ -f "coverage.out" ]; then
    go tool cover -func=coverage.out | tail -1
    echo "Para ver relatório HTML: go tool cover -html=coverage.out -o coverage.html"
else
    echo "⚠️  Arquivo de cobertura não encontrado"
fi
echo

echo "🎉 Pipeline testada com sucesso!"
echo "Agora você pode fazer commit e push com segurança."
echo
echo "Comandos úteis:"
echo "  make test     - Executar testes"
echo "  make lint     - Executar linter" 
echo "  make build    - Buildar exemplos"
echo "  make coverage - Gerar relatório de cobertura"
echo "  make ci       - Executar todos os checks"