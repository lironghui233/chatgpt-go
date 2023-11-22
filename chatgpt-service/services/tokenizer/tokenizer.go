package tokenizer

import (
	"bytes"
	"chatgpt-service/pkg/config"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sashabaranov/go-openai"
)

type tokenInfo struct {
	Code  int    `json:"code"`
	Count int    `json:"num_tokens"`
	Msg   string `json:"msg"`
}

func GetTokenCount(message openai.ChatCompletionMessage, model string) (int, error) {
	cnf := config.GetConf()
	url := fmt.Sprintf("%s/tokenizer/%s", cnf.DependOnServices.Tokenizer.Address, model)
	info := tokenInfo{}
	if err := postJSON(url, &message, &info); err != nil {
		return 0, err
	}
	if info.Code != 200 {
		return 0, fmt.Errorf("%v", info.Msg)
	}
	return info.Count, nil
}

func postJSON(url string, requestData *openai.ChatCompletionMessage, responseData *tokenInfo) error {
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(responseData)
}
