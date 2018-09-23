package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

const cfgFileName string = "./supporterCfg.json"

func getCfg() (supporterCfg, error) {
	_, err := os.Stat(cfgFileName)
	if err != nil {
		// assume cfg does not exist
		log.Printf("failed to read from cfg file:%s\n", err.Error())
		return getCfgFromArgs()
	}
	log.Printf("extract configuration from cfg file:%s\n", cfgFileName)
	return getCfgFromFile()
}

func getCfgFromFile() (supporterCfg, error) {
	file, err := ioutil.ReadFile(cfgFileName)
	if err != nil {
		return supporterCfg{}, err
	}
	var cfg supporterCfg
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return supporterCfg{}, err
	}
	return cfg, nil
}

func getCfgFromArgs() (supporterCfg, error) {
	args := os.Args
	if len(args) < 5 {
		return supporterCfg{}, errors.New("not enough args; recruit page url, company name, userId, and password are required")
	}
	url := args[1]
	company := args[2]
	mail := args[3]
	password := args[4]
	signinMethod := ""
	if len(args) > 5 {
		signinMethod = args[5]
	}
	cfg := supporterCfg{url, company, mail, password, signinMethod}
	return cfg, nil
}
