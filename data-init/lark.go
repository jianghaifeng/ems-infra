package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type recordResponseItem struct {
	Fields map[string]string `json:"fields"`
}

type recordResponseData struct {
	Has_more   bool                 `json:"has_more"`
	Page_token string               `json:"page_token"`
	Total      int                  `json:"total"`
	Items      []recordResponseItem `json:"items"`
}

type recordResponseBody struct {
	Code int                `json:"code"`
	Msg  string             `json:"msg"`
	Data recordResponseData `json:"data"`
}

func readLark(bitable string, table string, pageToken string) recordResponseBody {
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/records",
		bitable, table)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	q := req.URL.Query()
	q.Add("page_token", pageToken)
	q.Add("page_size", "500")
	q.Add("field_names", `["repository", "branch", "agile Team"]`)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", "Bearer "+getCurToken().Value)
	retry := 1
	var res *http.Response

	for ; retry <= 3; retry++ {
		res, err = client.Do(req)
		if err == nil && res.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(time.Second)
	}
	if retry > 3 {
		log.Panicln("retry out...")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var respBody recordResponseBody
	if err := json.Unmarshal(body, &respBody); err != nil {
		panic(err)
	}
	return respBody
}
