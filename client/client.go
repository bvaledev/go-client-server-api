package client

import (
	"DesafioClientServer/entity"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Message struct {
	Message string `json:"message"`
}

func Main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	isErrorMessage := resp.StatusCode != 200
	if isErrorMessage {
		var msg Message
		err := json.Unmarshal(data, &msg)
		if err != nil {
			panic(err.Error())
		}
		println(msg.Message)
		return
	}

	var quotation entity.Quotation
	json.Unmarshal(data, &quotation)
	save(quotation, resp)
}

func save(quotation entity.Quotation, resp *http.Response) {
	file, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	_, err = file.Write([]byte(fmt.Sprintf("Dólar: %f\n", quotation.Bid)))
	if err != nil {
		panic(err.Error())
	}
	io.Copy(os.Stdout, resp.Body)
}
