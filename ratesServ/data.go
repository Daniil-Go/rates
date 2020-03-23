package main

import "sync"

type Cache struct {
	Rates map[string]float64 `json:"rates,omitempty"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	mu    sync.Mutex
}

type Response struct {
	Base  string  `json:"base,omitempty"`
	Quote string  `json:"quote,omitempty"`
	Rate  float64 `json:"rate,omitempty"`
	Error string  `json:"error,omitempty"`
}

type SearchParams struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}
