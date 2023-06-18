package server

import (
	"DesafioClientServer/entity"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type QuotationAPIResponse struct {
	Usdbrl struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func Main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctxREQ, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()
	request, err := http.NewRequestWithContext(ctxREQ, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		errorResponse("", err, w)
		return
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		errorResponse("Request", err, w)
		return
	}
	defer response.Body.Close()

	var data QuotationAPIResponse
	json.NewDecoder(response.Body).Decode(&data)
	value, err := strconv.ParseFloat(data.Usdbrl.Bid, 64)
	if err != nil {
		errorResponse("", err, w)
		return
	}

	ctxDB, cancel := context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer cancel()

	quote := entity.NewQuotation(value)
	err = saveQuotation(ctxDB, quote)
	if err != nil {
		errorResponse("DB", err, w)
		return
	}

	json.NewEncoder(w).Encode(quote)
}

func connectDb() *sql.DB {
	db, err := sql.Open("sqlite3", "file:locked.sqlite?cache=shared")
	if err != nil {
		panic(err.Error())
	}
	tableQuotations := "CREATE TABLE IF NOT EXISTS quotations (id VARCHAR, bid REAL);"
	_, err = db.Exec(tableQuotations)
	if err != nil {
		panic(err.Error())
	}

	return db
}

func saveQuotation(ctx context.Context, quote *entity.Quotation) error {
	db := connectDb()
	defer db.Close()
	stmt, err := db.Prepare("insert into quotations(id, bid) values (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, quote.ID, quote.Bid)
	if err != nil {
		return err
	}
	return nil
}

func loadQuotations() {
	db := connectDb()
	defer db.Close()
	res, err := db.Query("select * from quotations")
	if err != nil {
		panic(err.Error())
	}
	var quotations []entity.Quotation
	for res.Next() {
		var q entity.Quotation
		res.Scan(&q.ID, &q.Bid)
		quotations = append(quotations, q)
	}
}

func errorResponse(name string, err error, w http.ResponseWriter) {
	if strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
		w.WriteHeader(http.StatusRequestTimeout)
		message := "{\"message\":\"" + name + " Timeout\"}"
		w.Write([]byte(message))
		log.Printf("[SERVER] %s", message)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	message := "{\"message\":\"Internal server error\"}"
	w.Write([]byte(message))
	log.Printf("[SERVER] %s", message)
}
