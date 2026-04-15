# Smart Ticket Triage

Sistema de classificação automática de chamados utilizando arquitetura orientada a eventos com Rails, Go, Redis e Gemini AI.

## Visão Geral

Este projeto demonstra uma arquitetura distribuída simples onde a criação de um ticket dispara um fluxo assíncrono de processamento e enriquecimento com inteligência artificial.

O objetivo é classificar automaticamente chamados técnicos, sugerindo categoria, prioridade e um resumo conciso.

## Arquitetura

* Rails responsável pelo CRUD de tickets
* Redis utilizado como barramento de mensagens via Pub Sub
* Go atuando como worker concorrente
* Gemini AI responsável pela análise e classificação dos tickets

Fluxo:

1. Usuário cria um ticket no Rails
2. Rails publica o evento no Redis
3. Worker em Go consome o evento
4. Worker chama a API do Gemini
5. Ticket é analisado e enriquecido

## Tecnologias

* Ruby on Rails
* Go (Golang)
* Redis
* Gemini API
* SQLite
* Docker

## Estrutura do Projeto

smart-ticket-triage/

* rails-app/   aplicação Rails (CRUD e publisher)
* go-worker/   worker em Go (consumer e processamento)

## Como rodar o projeto

### 1. Subir o Redis

```bash id="8yxr0q"
docker-compose up -d
```

---

### 2. Rodar o Rails

```bash id="x2c1ul"
cd rails-app
bundle install
rails db:migrate
rails s
```

Acesse:
http://localhost:3000/tickets

---

### 3. Rodar o worker em Go

```bash id="p4mqn9"
cd go-worker
go mod tidy
go run main.go
```

---

### 4. Configurar Gemini

Defina a variável de ambiente:

Windows PowerShell:

```powershell id="o0n5l9"
$env:GEMINI_API_KEY="SUA_CHAVE_AQUI"
```

## Exemplo de uso

Crie um ticket com qualquer título e descrição.

O sistema irá:

* enviar o ticket para o Redis
* processar o conteúdo no worker em Go
* utilizar IA para gerar uma classificação automática

A saída será exibida no console do worker, contendo:

* categoria sugerida
* prioridade estimada
* resumo do chamado

## Conceitos aplicados

* Arquitetura orientada a eventos
* Pub Sub com Redis
* Processamento assíncrono
* Concorrência com goroutines e channels
* Integração com IA generativa
* Separação de responsabilidades entre serviços

## Autor

Felício Melloni

LinkedIn:
https://www.linkedin.com/in/feliciomelloni/

GitHub:
https://github.com/fjgmelloni
