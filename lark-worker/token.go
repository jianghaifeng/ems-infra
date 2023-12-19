package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type TokenRequestBody struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type Token struct {
	Value  string `json:"tenant_access_token"`
	Expire int    `json:"expire"`
	time   time.Time
}

var curToken Token

func getCurToken() Token {
	return curToken
}

func refreshToken() error {
	if curToken.time.Add(time.Second * time.Duration(curToken.Expire-1000)).After(time.Now()) {
		return nil
	}

	reqBody := TokenRequestBody{
		AppId:     getConfigString("lark.auth.app_id"),
		AppSecret: getConfigString("lark.auth.app_secret"),
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", getConfigString("lark.auth.url"), strings.NewReader(string(payload)))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("status=" + res.Status)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, &curToken)
	curToken.time = time.Now()
	log.Println("new token acquired.")
	return nil
}
