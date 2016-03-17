package processing

import (
	"log"
	"io/ioutil"
	"encoding/json"
	"os"
	"sqlite"
	"strconv"
)

func readFile(file string) ([]HistoricalQuoteStruct,error) {
	data,err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("Error reading file ",file, " ",err)
		return nil,err
	}
	var hr HistoricalResponse
	err = json.Unmarshal(data,&hr)
	if err != nil {
		log.Println("Error unmarshalling file ",file, " ",err)
		return nil,err
	}
	return hr.Query.Results.Quotes,nil
}

func readDir(dir string) ([]string,error) {
	file, err := os.Open(dir)
	if err != nil {
		return nil,err
	}
	list, err := file.Readdirnames(20000)
	if err != nil {
		return nil,err
	}
	return list,nil
}

func update(partition string, in <- chan string) <-chan bool {
	out := make(chan bool)
	var dbname string
	switch partition {
	case "AC":
		dbname = "histoAC.db"
	case "DI":
		dbname = "histoDI.db"
	case "JP":
		dbname = "histoJP.db"
	case "QZ":
		dbname = "histoQZ.db"
	}
	var db sqlite.Sqlite
	db.Init(dbname)
	defer db.Destroy()

	go func() {
		for file := range in {
			hql,err := readFile(file)
			if err != nil {
				for _,quote := range hql {
					o, _ := strconv.ParseFloat(quote.Open, 64)
					c, _ := strconv.ParseFloat(quote.Close, 64)
					h, _ := strconv.ParseFloat(quote.High, 64)
					l, _ := strconv.ParseFloat(quote.Low, 64)
					v, _ := strconv.ParseInt(quote.Volume,10, 64)
					a, _ := strconv.ParseFloat(quote.AdjClose,64)
					record := &sqlite.DbHistTable{Date: quote.Date, Open: o,
						Close: c, High: h, Low: l,
						Volume: v, AdjClose: a}
					//db.CreateHistTable(quote.Symbol + "daily")
					db.InsertHist(record)
				}
			} else {
				log.Println("update failed for ",file, " ",err)
			}
		}
		out <- true
		close(out)
	}()
	return out
}

func splitter(datadir string) (<-chan string,<-chan string,<-chan string,<-chan string) {

	ac := make(chan string)
	di := make(chan string)
	jp := make(chan string)
	qz := make(chan string)

	filelist, err := readDir(datadir)
	if err == nil {
		close(ac)
		close(di)
		close(jp)
		close(qz)
		return ac,di,jp,qz
	}
	go func() {
		for _, file := range filelist {
			switch {
			case file[0] >= 'A' || file[0] <= 'C':
					ac <- file
			case file[0] >= 'D' || file[0] <= 'I':
					di <- file
			case file[0] >= 'J' || file[0] <= 'P' :
					jp <- file
			case file[0] >= 'Q' || file[0] <= 'Z':
					qz <- file
			}
		}
		close(ac)
		close(di)
		close(jp)
		close(qz)
	}()
	return ac,di,jp,qz
}


func run(datadir string) {
	ac,di,jp,qz := splitter(datadir)
	c1 := update("AC",ac)
	c2 := update("DI",di)
	c3 := update("JP",jp)
	c4 := update("QZ",qz)
	<-c1
	<-c2
	<-c3
	<-c4
}
