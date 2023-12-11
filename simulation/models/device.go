package models

type Device struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}
