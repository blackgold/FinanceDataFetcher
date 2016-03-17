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

func Historical(tasklist *[]*Task, cfg *config.Config, client *http.Client) []*Task {
	var errorTasks []*Task
	rate := time.Minute / time.Duration(cfg.Qps)
	throttle := time.Tick(rate)

	// sequential because we are anyways capped at 2k qps per hour
	for _, task := range *tasklist {
		if task.Retry > 3 {
			log.Println("FINAL FAIL : " + task.Symbol + " " + task.Start + "  " + task.End)
			continue
		}
		var baseurl string = "https://query.yahooapis.com/v1/public/yql?q="
		var query string = "select%20*%20from%20yahoo.finance.historicaldata%20where%20symbol%20%3D%20%22" + task.Symbol + "%22%20and%20startDate%20%3D%20%22" + task.Start + "%22%20and%20endDate%20%3D%20%22" + task.End + "%22&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="
		resp, err := client.Get(baseurl + query)
		if err != nil {
			log.Println("Error : "+task.Symbol+" "+task.Start+"  "+task.End+" ", err)
			task.Retry += 1
			tmp := *task
			errorTasks = append(errorTasks, &tmp)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error : "+task.Symbol+" "+task.Start+"  "+task.End+" ", err)
			task.Retry += 1
			tmp := *task
			errorTasks = append(errorTasks, &tmp)
			continue
		}
		var response HistoricalResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Error : "+task.Symbol+" "+task.Start+"  "+task.End+" ", err)
			task.Retry += 1
			tmp := *task
			errorTasks = append(errorTasks, &tmp)
			continue
		}
		if response.Query.Count > 0 {
			file := "data/" + task.Symbol + "-" + task.Start
			f, err := os.Create(file)
			if err == nil {
				_, err := f.Write(body)
				if err != nil {
					log.Println("Error : "+task.Symbol+" "+task.Start+"  "+task.End+" ", err)
					task.Retry += 1
					tmp := *task
					errorTasks = append(errorTasks, &tmp)
				}
			} else {
				log.Println("Error : "+task.Symbol+" "+task.Start+"  "+task.End+" ", err)
				task.Retry += 1
				tmp := *task
				errorTasks = append(errorTasks, &tmp)
			}
			if f != nil {
				f.Close()
			}
		}
		<-throttle
		log.Println("Done : ", task.Symbol)
	}
	return errorTasks
}

func RunHistorical(cfg *config.Config, client *http.Client) {
	symbols := readSymbols(cfg.SymbolsFile)
	var tasklist []*Task
	for _, symbol := range symbols {
		for i := 0; i < len(cfg.StartDates); i++ {
			tasklist = append(tasklist, &Task{Symbol: symbol, Start: cfg.StartDates[i], End: cfg.EndDates[i], Retry: 0})
		}
	}
	for len(tasklist) > 0 {
		tasklist = Historical(&tasklist, cfg, client)
	}
}

// cannot exploit parallelism here because rate-limiting is per ip
func HttpGet(query string, cfg *config.Config, client *http.Client) []string {
	var successquotes []string
	resp, err := client.Get(query)
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			var response DailyResponse
			err = json.Unmarshal(body, &response)
			if err == nil {
				if response.Query.Count > 0 {
					var db sqlite.Sqlite
					db.Init(cfg.DbFileName)
					defer db.Destroy()
					for _, quote := range response.Query.Results.Quotes {
						ask, _ := strconv.ParseFloat(quote.Ask, 64)
						averageDailyVolume, _ := strconv.ParseInt(quote.AverageDailyVolume, 10, 64)
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
						successquotes = append(successquotes, quote.Symbol)
						db.CreateDailyTable(quote.Symbol + "daily")
						db.InsertDaily(record)
						log.Println("Done : ", quote.Symbol)
					}
				} else {
					log.Println("Error response count is less than one ")
				}
			} else {
				log.Println("Error unmarshalling json body failed  ", err)
			}
		} else {
			log.Println("Error reading response body failed ", err)
		}
	} else {
		log.Println("Error http get failed ", err)
	}
	return successquotes
}

func Daily(tasklist *[]*Task, cfg *config.Config, client *http.Client) []*Task {
	var errorTasks []*Task
	var baseurl string = "https://query.yahooapis.com/v1/public/yql?q="
	var prequery string = "select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in%20"
	var postquery string = "&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="

	for i := 0; i < len(*tasklist); i += 5 {
		var selection string = "(%22"
		var limit = 0
		if i+5 < len(*tasklist) {
			limit = 5
		} else {
			limit = len(*tasklist) - i
		}
		for j := i; j < i+limit; j++ {
			if (*tasklist)[j].Retry > 4 {
				continue
			}
			if j == i+limit-1 {
				selection += (*tasklist)[j].Symbol + "%22)"
			} else {
				selection += (*tasklist)[j].Symbol + "%22%2C%22"
			}
		}
		if limit == 1 {
			selection = "(%22" + (*tasklist)[i].Symbol + "%22%2C%22" + "FAKESYMBOL" + "%22)"
		}
		//log.Println(selection)
		response := HttpGet(baseurl+prequery+selection+postquery, cfg, client)
		for k := i; k < i+limit; k++ {
			found := false
			for j := 0; j < len(response); j++ {
				if (*tasklist)[k].Symbol == response[j] {
					found = true
					break
				}
			}
			if !found {
				if (*tasklist)[k].Retry > 3 {
					log.Println("FINAL FAIL ", (*tasklist)[k].Symbol)
					continue
				}
				(*tasklist)[k].Retry += 1
				tmp := *(*tasklist)[k]
				errorTasks = append(errorTasks, &tmp)
			}
		}
	}
	return errorTasks
}

//TODO : implement using piepiline pattern
func RunDaily(cfg *config.Config, client *http.Client) {
	symbols := readSymbols(cfg.SymbolsFile)
	var tasklist []*Task
	for _, symbol := range symbols {
		tasklist = append(tasklist, &Task{Symbol: symbol, Start: "", End: "", Retry: 0})
	}
	i := 3
	for len(tasklist) > 0 {
		tasklist = Daily(&tasklist, cfg, client)
		time.Sleep(time.Duration(i*i)*time.Minute)
		i += 1
	}
}
