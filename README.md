![Arquitetura do Sistema](assets/fluxo.png)


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
* assets/      imagens e recursos

## Como rodar o projeto

### 1. Subir o Redis

```bash
docker-compose up -d
```

---

### 2. Rodar o Rails

```bash
cd rails-app
bundle install
rails db:migrate
rails s
```

Acesse:
http://localhost:3000/tickets

---

### 3. Rodar o worker em Go

```bash
cd go-worker
go mod tidy
go run main.go
```

---

### 4. Configurar Gemini

Defina a variável de ambiente:

Windows PowerShell:

```powershell
$env:GEMINI_API_KEY="SUA_CHAVE_AQUI"
```

---

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

---

## Funcionamento Interno

### Publicação de eventos no Rails

Ao criar um ticket, a aplicação Rails dispara automaticamente um evento para o Redis.

```ruby
class Ticket < ApplicationRecord
  after_create :send_to_queue

  def send_to_queue
    RedisPublisher.publish_ticket(self)
  end
end
```

---

### Service de publicação (RedisPublisher)

```ruby
class RedisPublisher
  def self.publish_ticket(ticket)
    redis = Redis.new

    payload = {
      id: ticket.id,
      title: ticket.title,
      description: ticket.description
    }

    redis.publish("tickets", payload.to_json)
  end
end
```

Responsável por:

* conexão com Redis
* serialização do ticket
* publicação no canal `tickets`

---

### Consumo e processamento no Go

O worker em Go consome os eventos e processa de forma assíncrona.

Fluxo do `main.go`:

1. Conecta ao Redis
2. Se inscreve no canal `tickets`
3. Recebe mensagens publicadas pelo Rails
4. Converte o JSON para struct
5. Envia para um channel interno
6. Processa com goroutine

---

### Concorrência com goroutines e channels

```go
jobs := make(chan TicketPayload, 100)

go worker(ctx, jobs)

for msg := range pubsub.Channel() {
    var payload TicketPayload
    json.Unmarshal([]byte(msg.Payload), &payload)
    jobs <- payload
}
```

Benefícios:

* processamento paralelo
* desacoplamento entre leitura e execução
* maior eficiência

---

### Integração com IA (Gemini)

O worker envia o ticket para o Gemini, que retorna:

* categoria
* prioridade
* resumo

A resposta é tratada como JSON e exibida no console.

---

## Conceitos aplicados

* Arquitetura orientada a eventos
* Pub Sub com Redis
* Processamento assíncrono
* Concorrência com goroutines e channels
* Integração com IA generativa
* Separação de responsabilidades

---

## Autor

Felício Melloni

LinkedIn:
https://www.linkedin.com/in/feliciomelloni/

GitHub:
https://github.com/fjgmelloni
