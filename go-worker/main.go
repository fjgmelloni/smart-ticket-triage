package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/genai"
)

type TicketPayload struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TicketAnalysis struct {
	Category string `json:"category"`
	Priority string `json:"priority"`
	Summary  string `json:"summary"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("erro ao conectar no redis: %v", err)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("variável GEMINI_API_KEY não encontrada")
	}

	aiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("erro ao criar client Gemini: %v", err)
	}

	log.Println("conectado no Redis")
	log.Println("cliente Gemini criado")

	pubsub := redisClient.Subscribe(ctx, "tickets")
	defer pubsub.Close()

	if _, err := pubsub.Receive(ctx); err != nil {
		log.Fatalf("erro ao assinar canal tickets: %v", err)
	}

	log.Println("inscrito no canal tickets")

	jobs := make(chan TicketPayload, 100)

	go worker(ctx, jobs, aiClient)

	go func() {
		ch := pubsub.Channel()

		for msg := range ch {
			var payload TicketPayload

			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				log.Printf("erro ao converter payload: %v", err)
				continue
			}

			log.Printf("ticket recebido do Redis: %+v\n", payload)
			jobs <- payload
		}
	}()

	waitForShutdown()
	log.Println("encerrando worker...")
}

func worker(ctx context.Context, jobs <-chan TicketPayload, aiClient *genai.Client) {
	for {
		select {
		case <-ctx.Done():
			log.Println("worker finalizado")
			return
		case ticket := <-jobs:
			processTicket(ctx, aiClient, ticket)
		}
	}
}

func processTicket(ctx context.Context, aiClient *genai.Client, ticket TicketPayload) {
	log.Println("processando ticket...")
	time.Sleep(1 * time.Second)

	analysis, err := analyzeTicket(ctx, aiClient, ticket)
	if err != nil {
		log.Printf("erro ao analisar ticket %d: %v", ticket.ID, err)
		return
	}

	fmt.Println("----------- TICKET ANALISADO -----------")
	fmt.Printf("ID: %d\n", ticket.ID)
	fmt.Printf("Título: %s\n", ticket.Title)
	fmt.Printf("Descrição: %s\n", ticket.Description)
	fmt.Printf("Categoria: %s\n", analysis.Category)
	fmt.Printf("Prioridade: %s\n", analysis.Priority)
	fmt.Printf("Resumo: %s\n", analysis.Summary)
	fmt.Println("----------------------------------------")
}

func analyzeTicket(ctx context.Context, aiClient *genai.Client, ticket TicketPayload) (TicketAnalysis, error) {
	prompt := fmt.Sprintf(`
Você é um classificador de chamados técnicos.

Analise o ticket abaixo e responda SOMENTE em JSON válido, sem markdown, sem explicações extras.

Formato esperado:
{
  "category": "string",
  "priority": "baixa|media|alta",
  "summary": "string"
}

Regras:
- category deve ser curta, como: financeiro, infraestrutura, autenticação, api, banco_de_dados, frontend, integrações
- priority deve ser exatamente: baixa, media ou alta
- summary deve ter no máximo 25 palavras

Ticket:
Título: %s
Descrição: %s
`, ticket.Title, ticket.Description)

	resp, err := aiClient.Models.GenerateContent(
		ctx,
		"gemini-3-flash-preview",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return TicketAnalysis{}, err
	}

	text := resp.Text()
	if text == "" {
		return TicketAnalysis{}, fmt.Errorf("resposta vazia da Gemini")
	}

	var analysis TicketAnalysis
	if err := json.Unmarshal([]byte(text), &analysis); err != nil {
		return TicketAnalysis{}, fmt.Errorf("erro ao converter resposta em JSON. resposta: %s | erro: %w", text, err)
	}

	return analysis, nil
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}