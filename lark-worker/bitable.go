package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type CopyTableRequestBody struct {
	Name           string `json:"name"`
	FolderToken    string `json:"folder_token"`
	WithoutContent bool   `json:"without_content"`
}

type TableOperationResponseBody struct {
	Data map[string]interface{} `json:"data"`
}

func refreshCurBiTable() {
	y, w := time.Now().ISOWeek()
	yw := fmt.Sprintf("%d-w%d", y, w)

	if getStatus().BitableName == yw {
		return
	}

	url := "https://open.feishu.cn/open-apis/bitable/v1/apps/" + getConfigString("lark.template.app_id") + "/copy"

	reqBody := CopyTableRequestBody{
		Name:           yw,
		FolderToken:    getConfigString("lark.template.folder_id"),
		WithoutContent: false,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))

	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+curToken.Value)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != http.StatusOK {
		panic(errors.New("status=" + res.Status))
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	res.Body.Close()
	var respBody TableOperationResponseBody
	if err := json.Unmarshal(body, &respBody); err != nil {
		panic(err)
	}

	app := respBody.Data["app"].(map[string]interface{})
	tableId := findCurTable(app["app_token"].(string))
	setBiTable(yw, app["app_token"].(string), tableId)
	log.Println("new bitable initialized.")
}

func findCurTable(bitable string) string {
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables", bitable)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

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
	var respBody TableOperationResponseBody
	if err := json.Unmarshal(body, &respBody); err != nil {
		panic(err)
	}
	app := respBody.Data["items"].([]interface{})
	item := app[0].(map[string]interface{})
	return item["table_id"].(string)
}

type Record struct {
	Item Item `json:"fields"`
}
type LarkRecordRequestBody struct {
	Records []Record `json:"records"`
}

func pushDataToLark(items []Item) {
	if len(items) == 0 {
		return
	}

	var body LarkRecordRequestBody
	for _, item := range items {
		body.Records = append(body.Records, Record{Item: item})
	}

	payload, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	url := "https://open.feishu.cn/open-apis/bitable/v1/apps/" + getStatus().BitableId +
		"/tables/" + getStatus().TableId + "/records/batch_create"
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+curToken.Value)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != http.StatusOK {
		panic(errors.New("status=" + res.Status))
	}

	setCursor(uint64(items[len(items)-1].Id))
	res.Body.Close()
}
