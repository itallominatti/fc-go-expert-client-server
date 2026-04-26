package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal("Erro ao criar requisição:", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Fatal("[TIMEOUT] Timeout de 300ms excedido ao aguardar resposta do servidor")
		}
		log.Fatal("Erro ao buscar cotação do servidor:", err)
	}
	defer resp.Body.Close()

	var cotacao CotacaoResponse
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		log.Fatal("Erro ao decodificar resposta JSON:", err)
	}

	content := fmt.Sprintf("Dólar: %s", cotacao.Bid)

	if err := os.WriteFile("cotacao.txt", []byte(content), 0644); err != nil {
		log.Fatal("Erro ao salvar cotacao.txt:", err)
	}

	log.Printf("Cotação salva com sucesso em cotacao.txt → %s\n", content)
}
