package models

type Document struct {
	Index map[string]string `json:"index"`
	Data  map[string]string `json:"data"`
}
