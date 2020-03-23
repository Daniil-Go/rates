package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {

	const (
		endpointURL = "https://api.exchangeratesapi.io/latest?base=USD"
		port        = "8080"
	)

	var rates Cache

	go func() {
		for {
			rates.getRates(endpointURL)
			log.Printf("Rates successfully updated")
			time.Sleep(60 * time.Second)
		}
	}()

	router := http.NewServeMux()

	router.HandleFunc("/rates", rates.ratesHandler)
	router.HandleFunc("/search", rates.searchHandler)

	log.Printf("Start server on port %v", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

	/*
		curl -v http://localhost:8080/rates
		http://localhost:8080/search?base=USD&quote=RUB

	*/
}

func (rates *Cache) getRates(reqURL string) {
	resp, err := http.Get(reqURL)
	if err != nil {
		log.Println(fmt.Errorf("%v during GET request", err))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(fmt.Errorf("%v during reading body of ratest", err))
		return
	}
	rates.mu.Lock()
	defer rates.mu.Unlock()
	if err = json.Unmarshal(body, &rates); err != nil {
		log.Println(fmt.Errorf("%v during Unmarshal", err))
		return
	}
}

func (rates *Cache) ratesHandler(wr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		sendError(wr, http.StatusBadRequest, fmt.Errorf("this is not GET request"))
		return
	}
	bytes, err := json.Marshal(rates.Rates)
	if err != nil {
		sendError(wr, http.StatusInternalServerError, fmt.Errorf("%v marshal response", err))
		return
	}
	_, err = wr.Write(bytes)
	if err != nil {
		sendError(wr, http.StatusInternalServerError, fmt.Errorf("%v write to response", err))
		return
	}

}

func (rates *Cache) searchHandler(wr http.ResponseWriter, req *http.Request) {
	searchParams := new(SearchParams)
	searchParams.Base = req.URL.Query().Get("base")
	searchParams.Quote = req.URL.Query().Get("quote")

	out := new(Response)
	if strings.EqualFold(searchParams.Base, rates.Base) {
		out.Base = searchParams.Base
	} else {
		fmt.Fprintf(wr, "The base \"%s\" is not correct", searchParams.Base)
		return
	}

	rates.mu.Lock()
	defer rates.mu.Unlock()
	rate, ok := rates.Rates[searchParams.Quote]
	if !ok {
		fmt.Fprintf(wr, "The quote \"%s\" is not correct", searchParams.Quote)
		return
	}

	out.Quote = searchParams.Quote
	out.Rate = rate

	fmt.Fprintf(wr, "Base is \"%s\" \nQuote is \"%s\" \nRate is\"%v\"",
		out.Base, out.Quote, out.Rate)
}

func sendError(wr http.ResponseWriter, statusCode int, err error) {
	wr.WriteHeader(statusCode)

	response := Response{
		Error: err.Error(),
	}

	bytes, err := json.Marshal(&response)
	if err != nil {
		log.Print(fmt.Errorf("%v during marshal response", err))
		return
	}

	_, err = wr.Write(bytes)
	if err != nil {
		log.Print(fmt.Errorf("%v during write response", err))
		return
	}
}
