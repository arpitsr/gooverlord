package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"com.ak.gooverlord/models"
)

func TestMain(m *testing.M) {
	ts := time.Now().Local().Format(time.RFC3339)
	log.Printf("%s", ts)
	url := "http://192.168.49.2:30726/logs"
	logEntry := models.LogEntry{
		Appname:   "testapp",
		Hostname:  "localhost",
		Message:   "Testing message",
		ID:        "1",
		Timestamp: ts,
	}
	b, err := json.Marshal(logEntry)
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

	body, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		fmt.Println(error)
	}

	fmt.Printf("Response Status:%+v\n", string(body))
	os.Exit(m.Run())
}

func TestQueryHandler(t *testing.T) {
	url := "http://192.168.49.2:30726/query"
	reqString := `{"q":"up.de","index":"17-10"}`
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

	// app := fiber.New()
	// app.Post("/query", handlers.Query)
	// resp, _ := app.Test(req)
	// if status := resp.StatusCode; status != http.StatusOK {
	// 	t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	// }
	expected := `[{"hits":[{"appname":"meln1ks","hostname":"up.de","message":"Pretty pretty pretty good","msgid":"ID141","timestamp":"2023-03-11T17:10:44.223Z"}],"estimatedTotalHits":1,"limit":20,"processingTimeMs":0,"query":"up.de"},{"hits":[],"limit":20,"processingTimeMs":0,"query":"up.de"}]`
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	if bodyString != expected {
		log.Fatalf("Received: %s vs Expected: %s", bodyString, expected)
	}
}
