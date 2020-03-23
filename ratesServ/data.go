package main

import "sync"

type Cache struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	mu    sync.Mutex
}

type Response struct {
	BaseName  string
	BaseRate  float64
	QuoteName string
	QuoteRate float64
	Rate      float64
	Error     string
}

type SearchParams struct {
	Base  string
	Quote string
}
