package processing

import (
	"bufio"
	"config"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sqlite"
	"strconv"
	"time"
)

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

func RunHistorical(cfg *config.Config) {

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
				log.Println("Error : "+symbol+" "+cfg.StartDates[i]+"  "+cfg.EndDates[i]+" ", err)
				continue
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error : "+symbol+" "+cfg.StartDates[i]+"  "+cfg.EndDates[i]+" ", err)
				continue
			}
			var response HistoricalResponse
			err = json.Unmarshal(body, &response)
			if err != nil {
				log.Println("Error : "+symbol+" "+cfg.StartDates[i]+"  "+cfg.EndDates[i]+" ", err)
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
					db.InsertHist(record)
				}
			} else {
				log.Println("Error : "+symbol+" "+cfg.StartDates[i]+"  "+cfg.EndDates[i]+" ", response.Query.Count)
				break
			}
			<-throttle
		}
		log.Println("Done : ", symbol)
	}
}

func RunDaily(cfg *config.Config) {

	symbol := readSymbols(cfg.SymbolsFile)
	rate := time.Minute / time.Duration(cfg.Qps)
	throttle := time.Tick(rate)

	for i := 0; i < len(symbol); i += 5 {
		var selection string = "(%22"
		var selections string
		var limit = 0
                if i+5 < len(symbol) {
                        limit = 5
                } else {
                        limit = len(symbol) - i
                }
		for j := i; j < i+limit ; j++ {
			if j == i+limit-1 {
				selection += symbol[j] + "%22)"
			} else {
				selection += symbol[j] + "%22%2C%22"
			}
 			selections += symbol[j] + ","
		}
		var baseurl string = "https://query.yahooapis.com/v1/public/yql?q="
		var query string = "select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in%20" + selection + "&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="
		resp, err := http.Get(baseurl + query)
		if err != nil {
			log.Println("Error : "+symbol[i]+" "+cfg.StartDates[i]+"  "+cfg.EndDates[i]+" ", err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error : "+symbol[i]+" "+cfg.StartDates[i]+"  "+cfg.EndDates[i]+" ", err)
			continue
		}
		log.Println(string(body[:]))
		var response DailyResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Error : "+symbol[i]+" "+cfg.StartDates[i]+"  "+cfg.EndDates[i]+" ", err)
			continue
		}
		if response.Query.Count > 0 {
			for _, quote := range response.Query.Results.Quotes {
				//var db sqlite.Sqlite
				//db.Init(cfg.DbFileName)
				//defer db.Destroy()
				//db.CreateDailyTable(quote.Symbol + "daily")
				ask, _ := strconv.ParseFloat(quote.Ask, 64)
				averageDailyVolume, _ := strconv.ParseInt(quote.AverageDailyVolume,10, 64)
				bid, _ := strconv.ParseFloat(quote.Bid, 64)
				bookValue, _ := strconv.ParseFloat(quote.BookValue, 64)
				earningsShare, _ := strconv.ParseFloat(quote.EarningsShare, 64)
				daysLow, _ := strconv.ParseFloat(quote.DaysLow, 64)
				daysHigh, _ := strconv.ParseFloat(quote.DaysHigh, 64)
				yearLow, _ := strconv.ParseFloat(quote.YearLow, 64)
				yearHigh, _ := strconv.ParseFloat(quote.ChangeFromYearLow, 64)
				changeFromYearLow, _ := strconv.ParseFloat(quote.ChangeFromYearLow, 64)
				changeFromYearHigh, _ := strconv.ParseFloat(quote.ChangeFromYearHigh, 64)
				fiftydayMovingAverage, _ := strconv.ParseFloat(quote.FiftydayMovingAverage, 64)
				pEGRatio, _ := strconv.ParseFloat(quote.PEGRatio, 64)
				shortRatio, _ := strconv.ParseFloat(quote.ShortRatio, 64)
				volume, _ := strconv.ParseInt(quote.Volume, 10, 64)
				dividendYield, _ := strconv.ParseFloat(quote.DividendYield, 64)
				twoHundreddayMovingAverage, _ := strconv.ParseFloat(quote.TwoHundreddayMovingAverage, 64)
				changeFromTwoHundreddayMovingAverage, _ := strconv.ParseFloat(quote.ChangeFromTwoHundreddayMovingAverage, 64)
				changeFromFiftydayMovingAverage, _ := strconv.ParseFloat(quote.ChangeFromFiftydayMovingAverage, 64)
				t := time.Now().Local()
    				date := t.Format("2006-01-02")
				record := &sqlite.DbDailyTable{Date: date,
					Ask:                                         ask,
					AverageDailyVolume:                          averageDailyVolume,
					Bid:                                         bid,
					BookValue:                                   bookValue,
					EarningsShare:                               earningsShare,
					PEGRatio:                                    pEGRatio,
					ShortRatio:                                  shortRatio,
					Volume:                                      volume,
					DividendYield:                               dividendYield,
					ChangeinPercent:                             quote.ChangeinPercent,
					DaysLow:                                     daysLow,
					DaysHigh:                                    daysHigh,
					YearLow:                                     yearLow,
					YearHigh:                                    yearHigh,
					MarketCapitalization:                        quote.MarketCapitalization,
					EBITDA:                                      quote.EBITDA,
					ChangeFromYearLow:                           changeFromYearLow,
					PercentChangeFromYearLow:                    quote.PercentChangeFromYearLow,
					ChangeFromYearHigh:                          changeFromYearHigh,
					PercentChangeFromYearHigh:                   quote.PercentChangeFromYearHigh,
					FiftydayMovingAverage:                       fiftydayMovingAverage,
					TwoHundreddayMovingAverage:                  twoHundreddayMovingAverage,
					ChangeFromTwoHundreddayMovingAverage:        changeFromTwoHundreddayMovingAverage,
					PercentChangeFromTwoHundreddayMovingAverage: quote.PercentChangeFromTwoHundreddayMovingAverage,
					ChangeFromFiftydayMovingAverage:             changeFromFiftydayMovingAverage,
					PercentChangeFromFiftydayMovingAverage:      quote.PercentChangeFromFiftydayMovingAverage}
				//db.InsertDaily(record)
				log.Println(record)
				log.Println("Done : ", quote.Symbol)
			}
		} else {
			log.Println("Error : "+selections+" "+cfg.StartDates[i]+"  "+cfg.EndDates[i]+" ", response.Query.Count)
			break
		}
		<-throttle
	}

}
