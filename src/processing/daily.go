package processing

type DailyQuoteStruct struct {
	Ssymbol                                     string `json:"symbol"`
	Ask                                         string `json:"Ask"`
	AverageDailyVolume                          string `json:"AverageDailyVolume"`
	Bid                                         string `json:"Bid"`
	AskRealtime                                 string `json:"AskRealtime"`
	BidRealtime                                 string `json:"BidRealtime"`
	BookValue                                   string `json:"BookValue"`
	ChangePercentChange                         string `json:"Change_PercentChange"`
	Change                                      string `json:"Change"`
	Commission                                  string `json:"Commission"`
	Currency                                    string `json:"Currency"`
	ChangeRealtime                              string `json:"ChangeRealtime"`
	AfterHoursChangeRealtime                    string `json:"AfterHoursChangeRealtime"`
	DividendShare                               string `json:"DividendShare"`
	LastTradeDate                               string `json:"LastTradeDate"`
	TradeDate                                   string `json:"TradeDate"`
	EarningsShare                               string `json:"EarningsShare"`
	ErrorIndication                             string `json:"ErrorIndicationreturnedforsymbolchangedinvalid"`
	EPSEstimateCurrentYear                      string `json:"EPSEstimateCurrentYear"`
	EPSEstimateNextYear                         string `json:"EPSEstimateNextYear"`
	EPSEstimateNextQuarter                      string `json:"EPSEstimateNextQuarter"`
	DaysLow                                     string `json:"DaysLow"`
	DaysHigh                                    string `json:"DaysHigh"`
	YearLow                                     string `json:"YearLow"`
	YearHigh                                    string `json:"YearHigh"`
	HoldingsGainPercent                         string `json:"HoldingsGainPercent"`
	AnnualizedGain                              string `json:"AnnualizedGain"`
	HoldingsGain                                string `json:"HoldingsGain"`
	HoldingsGainPercentRealtime                 string `json:"HoldingsGainPercentRealtime"`
	HoldingsGainRealtime                        string `json:"HoldingsGainRealtime"`
	MoreInfo                                    string `json:"MoreInfo"`
	OrderBookRealtime                           string `json:"OrderBookRealtime"`
	MarketCapitalization                        string `json:"MarketCapitalization"`
	MarketCapRealtime                           string `json:"MarketCapRealtime"`
	EBITDA                                      string `json:"EBITDA"`
	ChangeFromYearLow                           string `json:"ChangeFromYearLow"`
	PercentChangeFromYearLow                    string `json:"PercentChangeFromYearLow"`
	LastTradeRealtimeWithTime                   string `json:"LastTradeRealtimeWithTime"`
	ChangePercentRealtime                       string `json:"ChangePercentRealtime"`
	ChangeFromYearHigh                          string `json:"ChangeFromYearHigh"`
	PercentChangeFromYearHigh                   string `json:"PercebtChangeFromYearHigh"`
	LastTradeWithTime                           string `json:"LastTradeWithTime"`
	LastTradePriceOnly                          string `json:"LastTradePriceOnly"`
	HighLimit                                   string `json:"HighLimit"`
	LowLimit                                    string `json:"LowLimit"`
	DaysRange                                   string `json:"DaysRange"`
	DaysRangeRealtime                           string `json:"DaysRangeRealtime"`
	FiftydayMovingAverage                       string `json:"FiftydayMovingAverage"`
	TwoHundreddayMovingAverage                  string `json:"TwoHundreddayMovingAverage"`
	ChangeFromTwoHundreddayMovingAverage        string `json:"ChangeFromTwoHundreddayMovingAverage"`
	Name                                        string `json:"Name"`
	Notes                                       string `json:"Notes"`
	Open                                        string `json:"Open"`
	PreviousClose                               string `json:"PreviousClose"`
	PricePaid                                   string `json:"PricePaid"`
	ChangeinPercent                             string `json:"ChangeinPercent"`
	PriceSales                                  string `json:"PriceSales"`
	PriceBook                                   string `json:"PriceBook"`
	ExDividendDate                              string `json:"ExDividendDate"`
	PERatio                                     string `json:"PERatio"`
	DividendPayDate                             string `json:"DividendPayDate"`
	PERatioRealtime                             string `json:"PERatioRealtime"`
	PEGRatio                                    string `json:"PEGRatio"`
	PriceEPSEstimateCurrentYear                 string `json:"PriceEPSEstimateCurrentYear"`
	PriceEPSEstimateNextYear                    string `json:"PriceEPSEstimateNextYear"`
	Symbol                                      string `json:"Symbol"`
	SharesOwned                                 string `json:"SharesOwned"`
	ShortRatio                                  string `json:"ShortRatio"`
	LastTradeTime                               string `json:"LastTradeTime"`
	TickerTrend                                 string `json:"TickerTrend"`
	OneyrTargetPrice                            string `json:"OneyrTargetPrice"`
	Volume                                      string `json:"Volume"`
	HoldingsValue                               string `json:"HoldingsValue"`
	HoldingsValueRealtime                       string `json:"HoldingsValueRealtime"`
	YearRange                                   string `json:"YearRange"`
	DaysValueChange                             string `json:"DaysValueChange"`
	DaysValueChangeRealtime                     string `json:"DaysValueChangeRealtime"`
	StockExchange                               string `json:"StockExchange"`
	DividendYield                               string `json:"DividendYield"`
	PercentChange                               string `json:"PercentChange"`
	PercentChangeFromTwoHundreddayMovingAverage string `json:"PercentChangeFromTwoHundreddayMovingAverage"`
	ChangeFromFiftydayMovingAverage             string `json:"ChangeFromFiftydayMovingAverage"`
	PercentChangeFromFiftydayMovingAverage      string `json:"PercentChangeFromFiftydayMovingAverage"`
}
type DailyResultsStruct struct {
	Quotes []DailyQuoteStruct `json:"quote"`
}

type DailyQueryStruct struct {
	Count   int                `json:"count"`
	Created string             `json:"created"`
	Lang    string             `json:"lang"`
	Results DailyResultsStruct `json:"results"`
}

type DailyResponse struct {
	Query DailyQueryStruct `json:"query"`
}


