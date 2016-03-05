package config

import (
        "encoding/json"
        "io/ioutil"
       )

type Config struct {
     StartDates      []string `json:"StartDates"`
     EndDates        []string `json:"EndDates"`
     SymbolsFile     string   `json:"SymbolsFile"`
     DbFileName      string   `json:"DbFileName"`	
     Qps	     int      `json:"Qps"`

}

func Parse(file string) (*Config, error) {
	cfg := &Config{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, cfg)
	return cfg, err
}
