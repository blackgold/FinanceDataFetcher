package main

import (
	"bufio"
	"os"
	"log"
//        "net/http"
//        "io/ioutil"
//	"encoding/json"
)


type QuoteStruct struct {
  Symbol   string `json:"Symbol"`
  Date     string `json:"Date"`
  Open     string `json:"Open"`
  High     string `json:"High"`
  Low      string `json:"Low"`
  Close    string `json:"Close"`
  Volume   string `json:"Volume"`
 Adj_Close string `json:"Adj_Close"`
}

type ResultsStruct struct {
  Quote   []QuoteStruct `json:"quote"`  
}

type QueryStruct struct{
  Count     int           `json:"count"`
  Created   string        `json:"created"`
  Lang      string        `json:"lang"`
  Results   ResultsStruct `json:"results"`
}

type Response struct {
  Query   QueryStruct `json:"query"`
}

func readSymbols(file string) []string {
    var symbols []string
    fd, err := os.Open(file)
    if err != nil {
        log.Fatal(err)
    }
    defer fd.Close()

    scanner := bufio.NewScanner(fd)
    for scanner.Scan() {
        symbols = append(symbols,scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return symbols
}

func main() {
 
 //start := [5]string{"2016-01-01","2015-01-01","2014-01-01","2013-01-01","2012-01-01"}
 //end :=  [5]string{"2016-02-29","2015-12-31","2014-12-31","2013-12-31","2012-12-31"}
 log.Print(readSymbols("symbols/symbols.txt")) 
}
 
/*
 var baseurl string = "https://query.yahooapis.com/v1/public/yql?q="
 var  query string = "select%20*%20from%20yahoo.finance.historicaldata%20where%20symbol%20%3D%20%22TEAM%22%20and%20startDate%20%3D%20%22" + start + "%22%20and%20endDate%20%3D%20%22" + end + "%22&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="
 resp, err := http.Get(baseurl + query)
 if err !=nil {
    log.Fatal(err)
 }
 defer resp.Body.Close()
 body, err := ioutil.ReadAll(resp.Body)
 if err != nil {
    log.Fatal(err)
 }
 var result Response
 err = json.Unmarshal(body, &result)
 if err != nil {
    log.Fatal(">",err)
 }
 log.Print(result)
*/
