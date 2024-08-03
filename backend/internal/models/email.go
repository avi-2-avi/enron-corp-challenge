package models

type Email struct {
    Date      string `json:"date"`
    From      string `json:"from"`
    To        string `json:"to"`
    Subject   string `json:"subject"`
    Content   string `json:"content"`
    Path      string `json:"path"`
    Timestamp string `json:"@timestamp"`
}
