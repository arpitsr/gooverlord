package fulltexttest_test

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

func seedData() []models.LogEntry {
	return []models.LogEntry{
		{
			Appname:   "mytestapp",
			Hostname:  "mytesthost",
			Message:   "cat says meow 1",
			ID:        "1",
			Timestamp: ts.Format(time.RFC3339),
		},
		{
			Appname:   "mytestapp",
			Hostname:  "mytesthost",
			Message:   "cat says meow 2",
			ID:        "2",
			Timestamp: ts.Format(time.RFC3339),
		},
		{
			Appname:   "mytestapp",
			Hostname:  "mytesthost",
			Message:   "cat says meow 3",
			ID:        "3",
			Timestamp: ts.Format(time.RFC3339),
		},
	}
}

func TestMain(m *testing.M) {
	url := fmt.Sprintf("%s/logs", SERVER_URL)

	b, err := json.Marshal(seedData())
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

func TestFullTextSearch(t *testing.T) {
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("%s/query", SERVER_URL)
	hourMinute := ts.Format(config.DATE_INDEX_FORMAT)
	fmt.Print(hourMinute)
	// set up test data
	q := models.Query{
		Index:       hourMinute,
		SearchQuery: "cat",
	}
	b, err := json.Marshal(q)
	if err != nil {
		t.Fatalf("Could not marshal the req %e", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
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

	var res *meilisearch.SearchResponse

	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		t.Fatalf("Could not decode response, %e", err)
	}

	// check that the returned value is not nil
	if res == nil {
		t.Error("Expected a non-nil search response")
	}

	// check that the returned value has expected properties
	if res.EstimatedTotalHits != 3 {
		t.Errorf("Expected %d estimated total hits, but got %d", 3, res.EstimatedTotalHits)
	}

	// check that the returned value contains the expected hits
	expectedHits := `{"hits":[{"appname":"mytestapp","hostname":"mytesthost","id":"1","message":"cat says meow 1","timestamp":"2017-05-11T17:59:16+05:30"},{"appname":"mytestapp","hostname":"mytesthost","id":"2","message":"cat says meow 2","timestamp":"2017-05-11T17:59:16+05:30"},{"appname":"mytestapp","hostname":"mytesthost","id":"3","message":"cat says meow 3","timestamp":"2017-05-11T17:59:16+05:30"}],"estimatedTotalHits":3,"limit":20,"processingTimeMs":0,"query":"cat"}`
	var expectedStruct meilisearch.SearchResponse

	err = json.Unmarshal([]byte(string(expectedHits)), &expectedStruct)
	if err != nil {
		log.Printf("Could not unmarshal expected Hits %e", err)
	}

	if !reflect.DeepEqual(expectedStruct, *res) {
		t.Errorf("\n%+v\n%+v\n", expectedStruct, *res)
	}
}
