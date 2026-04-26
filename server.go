package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type AwesomeAPIResponse struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

var db *sql.DB

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./cotacoes.db"
	}

	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("Erro ao abrir banco de dados:", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		bid        TEXT    NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Erro ao criar tabela:", err)
	}

	http.HandleFunc("/cotacao", handleCotacao)
	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCotacao(w http.ResponseWriter, r *http.Request) {
	apiCtx, apiCancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer apiCancel()

	cotacao, err := fetchCotacao(apiCtx)
	if err != nil {
		if apiCtx.Err() == context.DeadlineExceeded {
			log.Println("[TIMEOUT] Timeout ao buscar cotação na API externa:", err)
			http.Error(w, "Timeout ao buscar cotação na API externa", http.StatusGatewayTimeout)
			return
		}
		log.Println("[ERRO] Erro ao buscar cotação na API externa:", err)
		http.Error(w, "Erro interno ao buscar cotação", http.StatusInternalServerError)
		return
	}

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer dbCancel()

	if err := saveCotacao(dbCtx, cotacao.USDBRL.Bid); err != nil {
		if dbCtx.Err() == context.DeadlineExceeded {
			log.Println("[TIMEOUT] Timeout ao persistir cotação no banco de dados:", err)
		} else {
			log.Println("[ERRO] Erro ao persistir cotação no banco de dados:", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CotacaoResponse{Bid: cotacao.USDBRL.Bid})
}

func fetchCotacao(ctx context.Context) (*AwesomeAPIResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result AwesomeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func saveCotacao(ctx context.Context, bid string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO cotacoes (bid) VALUES (?)", bid)
	return err
}
