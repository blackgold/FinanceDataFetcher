package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
	"log"
)

type DbHistTable struct {
	Date     string
	Open     float64
	Close    float64
	AdjClose float64
	High     float64
	Low      float64
	Volume   int64
}

type DbDailyTable struct {
	Date                                        string
	Ask                                         float64
	AverageDailyVolume                          int64
	Bid                                         float64
	BookValue                                   float64
	EarningsShare                               float64
	DaysLow                                     float64
	DaysHigh                                    float64
	YearLow                                     float64
	YearHigh                                    float64
	MarketCapitalization                        string
	EBITDA                                      string
	ChangeFromYearLow                           float64
	PercentChangeFromYearLow                    string
	ChangeFromYearHigh                          float64
	PercentChangeFromYearHigh                   string
	LastTradePriceOnly                          float64
	FiftydayMovingAverage                       float64
	TwoHundreddayMovingAverage                  float64
	ChangeFromTwoHundreddayMovingAverage        float64
	PercentChangeFromTwoHundreddayMovingAverage string
	ChangeFromFiftydayMovingAverage             float64
	PercentChangeFromFiftydayMovingAverage      string
	PEGRatio                                    float64
	ShortRatio                                  float64
	Volume                                      int64
	DividendYield                               float64
	ChangeinPercent                             string
}

type Sqlite struct {
	dbmap *gorp.DbMap
}

func (s *Sqlite) Init(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	db.Exec("PRAGMA synchronous=NORMAL")
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA temp_store=MEMORY")
	db.Exec("PRAGMA cache_size=4096")
	s.dbmap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	return nil
}

func (s *Sqlite) Destroy() {
	s.dbmap.Db.Close()
}

func (s *Sqlite) CreateHistTable(name string) error {
	tmp := s.dbmap.AddTableWithName(DbHistTable{}, name).SetKeys(false, "Date")
	tmp.ColMap("Date").SetMaxSize(10)
	err := s.dbmap.CreateTablesIfNotExists()
	if err != nil {
		s.dbmap.Db.Close()
		return err
	}
	return nil
}

func (s *Sqlite) CreateDailyTable(name string) error {
	tmp := s.dbmap.AddTableWithName(DbDailyTable{}, name).SetKeys(false, "Date")
	tmp.ColMap("Date").SetMaxSize(10)
	err := s.dbmap.CreateTablesIfNotExists()
	if err != nil {
		s.dbmap.Db.Close()
		return err
	}
	return nil
}

func (s *Sqlite) InsertHist(value *DbHistTable) error {
	err := s.dbmap.Insert(value)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *Sqlite) InsertDaily(value *DbDailyTable) error {
	err := s.dbmap.Insert(value)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
