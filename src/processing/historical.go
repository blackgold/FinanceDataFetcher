package processing

type HistoricalQuoteStruct struct {
	Symbol   string `json:"Symbol"`
	Date     string `json:"Date"`
	Open     string `json:"Open"`
	High     string `json:"High"`
	Low      string `json:"Low"`
	Close    string `json:"Close"`
	Volume   string `json:"Volume"`
	AdjClose string `json:"Adj_Close"`
}

type HistoricalResultsStruct struct {
	Quotes []HistoricalQuoteStruct `json:"quote"`
}

type HistoricalQueryStruct struct {
	Count   int                     `json:"count"`
	Created string                  `json:"created"`
	Lang    string                  `json:"lang"`
	Results HistoricalResultsStruct `json:"results"`
}

type HistoricalResponse struct {
	Query HistoricalQueryStruct `json:"query"`
}
