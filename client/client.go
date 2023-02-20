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

func Main() {
	ctx, _ := context.WithTimeout(context.Background(), 300*time.Millisecond)
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var quotation entity.Quotation
	json.Unmarshal(data, &quotation)

	file, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write([]byte(fmt.Sprintf("Dólar: %f\n", quotation.Bid)))
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, resp.Body)
}
