package historical

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sqlite"
	"config"
	"strconv"
	"time"
)

type QuoteStruct struct {
	Symbol   string `json:"Symbol"`
	Date     string `json:"Date"`
	Open     string `json:"Open"`
	High     string `json:"High"`
	Low      string `json:"Low"`
	Close    string `json:"Close"`
	Volume   string `json:"Volume"`
	AdjClose string `json:"Adj_Close"`
}

type ResultsStruct struct {
	Quotes []QuoteStruct `json:"quote"`
}

type QueryStruct struct {
	Count   int           `json:"count"`
	Created string        `json:"created"`
	Lang    string        `json:"lang"`
	Results ResultsStruct `json:"results"`
}

type Response struct {
	Query QueryStruct `json:"query"`
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
		symbols = append(symbols, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return symbols
}

func Run(cfg *config.Config) {

        symbols := readSymbols(cfg.SymbolsFile)
	rate := time.Minute / time.Duration(cfg.Qps)
	throttle := time.Tick(rate)
	for _, symbol := range symbols {
		var db sqlite.Sqlite
		db.Init(cfg.DbFileName)
		defer db.Destroy()
		db.CreateHistTable(symbol + "history")
		for i := 0; i < 10; i++ {
			var baseurl string = "https://query.yahooapis.com/v1/public/yql?q="
			var query string = "select%20*%20from%20yahoo.finance.historicaldata%20where%20symbol%20%3D%20%22" + symbol + "%22%20and%20startDate%20%3D%20%22" + cfg.StartDates[i] + "%22%20and%20endDate%20%3D%20%22" + cfg.EndDates[i] + "%22&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="
			resp, err := http.Get(baseurl + query)
			if err != nil {
				log.Println("Error : " + symbol + " " + cfg.StartDates[i] + "  " + cfg.EndDates[i] + " ", err)
				continue
			}
			defer resp.Body.Close()
                        
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error : " + symbol + " " + cfg.StartDates[i] + "  " + cfg.EndDates[i] + " ", err)
				continue
			}
			var response Response
			err = json.Unmarshal(body, &response)
			if err != nil {
				log.Println("Error : " + symbol + " " + cfg.StartDates[i] + "  " + cfg.EndDates[i] + " ", err)
				continue
			}
			if response.Query.Count > 0 {
				for _, quote := range response.Query.Results.Quotes {
					open, _ := strconv.ParseFloat(quote.Open, 64)
					lose, _ := strconv.ParseFloat(quote.Close, 64)
					high, _ := strconv.ParseFloat(quote.High, 64)
					low, _ := strconv.ParseFloat(quote.Low, 64)
					adj, _ := strconv.ParseFloat(quote.AdjClose, 64)
					vol, _ := strconv.ParseInt(quote.Volume, 10, 64)
					record := &sqlite.DbHistTable{Date: quote.Date, Open: open, Close: lose, High: high, Low: low, Volume: vol, AdjClose: adj}
					db.Insert(record)
				}
			} else {
				log.Println("Error : " + symbol + " " + cfg.StartDates[i] + "  " + cfg.EndDates[i] + " ", response.Query.Count)
				break
			}
			<-throttle
		}
		log.Println("Done : ", symbol)
	}
}
