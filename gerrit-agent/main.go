package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/segmentio/kafka-go"
)

func writeByConn(message string) {
	topic := "test"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9094", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte(message)},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

type Change struct {
	Id     string
	Number uint64 `json:"number"`
}

type MergedEvent struct {
	Change     Change `json:"change"`
	EventType  string `json:"type"`
	InstanceId string `json:"instanceId"`
}

func postEvent(c *gin.Context) {
	var event MergedEvent
	if err := c.BindJSON(&event); err != nil {
		log.Println("error binding json:", err)
		return
	}

	event.Change.Id = strconv.FormatUint(event.Change.Number, 10)

	log.Printf("receive %s event from %s, id = %s\n",
		event.EventType, event.InstanceId, event.Change.Id)

	writeByConn(event.EventType)

	if event.EventType == "change-merged" {
		changeIds.LoadOrStore(event.Change.Id, Doc{Status: "accepted", Time: time.Now()})
	}
}

func getChangeDetails(id string) ([]byte, error) {
	req, _ := http.NewRequest("GET", "http://gerrit.haifeng.com/a/changes/"+id+"/detail", nil)
	req.SetBasicAuth(getConfigString("gerrit.username"), getConfigString("gerrit.password"))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("getChangeDetails, status=" + resp.Status)
		return nil, errors.New("status=" + resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return bytes.TrimLeft(b, ")]}'"), nil
}

type Doc struct {
	Status string    `json:"status"`
	Time   time.Time `json:"@timestamp"`
}

type DetailDoc struct {
	Detail interface{} `json:"detail"`
	Doc
}

var changeIds sync.Map

func processDoc(key, value any) bool {
	id := key.(string)
	doc := value.(Doc)

	if time.Since(doc.Time).Minutes() > 1 {
		log.Println("timeout for id:", id, doc.Time.Format(time.DateTime))
		changeIds.Delete(id)
	}

	if doc.Status == "accepted" && getDoc(id) {
		log.Println("found doc in es, change status to complete, id=", id)
		doc.Status = "complete"
		changeIds.Store(id, doc)
	}

	if doc.Status == "fail" || doc.Status == "complete" {
		return true
	}

	body, err := getChangeDetails(id)
	if err != nil {
		if doc.Status == "accepted" {
			doc.Status = "retry"
		} else {
			doc.Status = "fail"
		}
	} else {
		doc.Status = "complete"
	}
	changeIds.Store(id, doc)

	if doc.Status == "retry" {
		return false
	}

	detailDoc := DetailDoc{Doc: doc}
	if err := json.Unmarshal(body, &detailDoc.Detail); err != nil {
		log.Println(err)
	}
	if err := addDoc(id, detailDoc); err != nil {
		log.Println(err)
	}
	if err := addCommit(body); err != nil {
		log.Println(err)
	}
	return true
}

func addDoc(id string, doc DetailDoc) error {
	body, _ := json.Marshal(doc)
	req, _ := http.NewRequest("POST", getConfigString("es.url")+"/_doc/"+id, bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(getConfigString("es.username"), getConfigString("es.password"))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.New("add Doc status=" + resp.Status)
	}
	return nil
}

func addCommit(body []byte) error {
	resp, err := http.Post(getConfigString("ems.url"), "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.New("add Commit status=" + resp.Status)
	}
	return nil
}

func getDoc(id string) bool {
	req, _ := http.NewRequest("GET", getConfigString("es.url")+"/_doc/"+id, nil)
	req.SetBasicAuth(getConfigString("es.username"), getConfigString("es.password"))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln("error get doc:", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return false
	}
	return true
}

func processDocs() {
	changeIds.Range(processDoc)
}

func process() {
	for {
		processDocs()
		time.Sleep(time.Second * 2)
	}
}

func main() {
	loadConfig()
	router := gin.Default()
	router.POST("/events", postEvent)

	go process()
	router.Run(":8080")
}
