package models

type LogEntry struct {
	Appname   string `json:"appname"`
	Hostname  string `json:"hostname"`
	Message   string `json:"message"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
}
