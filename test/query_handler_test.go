package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"com.ak.gooverlord/config"
	"com.ak.gooverlord/models"
	"github.com/meilisearch/meilisearch-go"
)

const SERVER_URL = "http://localhost:3000"

var ts = time.Unix(1494505756, 0)

func TestMain(m *testing.M) {
	url := fmt.Sprintf("%s/logs", SERVER_URL)
	logEntry := models.LogEntry{
		Appname:   "testapp",
		Hostname:  "localhost",
		Message:   "Testing message",
		ID:        "1",
		Timestamp: ts.Format(time.RFC3339),
	}
	b, err := json.Marshal([]models.LogEntry{logEntry})
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(b)))
	if err != nil {
		log.Fatalf(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		fmt.Println(error)
	}

	fmt.Printf("Response Status:%+v\n", string(body))
	os.Exit(m.Run())
}

func TestQueryHandler(t *testing.T) {
	url := fmt.Sprintf("%s/query", SERVER_URL)
	hourMinute := ts.Format(config.DATE_INDEX_FORMAT)
	reqString := fmt.Sprintf(`{"q":"Testing","index":"%s"}`, hourMinute)
	requestBody := []byte(reqString)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal("Error in send query req")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	var res []*meilisearch.SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		t.Fatalf("Could not decode response, %e", err)
	}

	var exp []*meilisearch.SearchResponse
	expected := `[{"hits":[{"appname":"testapp","hostname":"localhost","id":"1","message":"Testing message","timestamp":"2017-05-11T17:59:16+05:30"}],"estimatedTotalHits":1,"limit":20,"processingTimeMs":0,"query":"Testing"}]`

	err = json.Unmarshal([]byte(expected), &exp)
	if err != nil {
		t.Fatalf("Could not decode exp, %e", err)
	}
	if reflect.DeepEqual(exp, res) {
		log.Printf("Success")
	} else {
		t.Fatalf("Exp: %+v, Returned: %+v", exp, res)
	}
}
