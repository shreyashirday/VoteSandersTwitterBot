package main

import (
       "encoding/json"
       "os"
)

type Configuration struct {
     ConsumerKey    string
     ConsumerSecret string
     AccessToken    string
     AccessKey      string
     MetamindKey    string
}

func getConfig(fileName string) (Configuration, error) {
     file, _ := os.Open(fileName)
     decoder := json.NewDecoder(file)
     configuration := Configuration{}
     err := decoder.Decode(&configuration)
     return configuration, err
}