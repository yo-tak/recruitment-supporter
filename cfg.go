package main

type supporterCfg struct {
	URL          string `json:"url"`
	Company      string `json:"company"`
	Mail         string `json:"mail"`
	Password     string `json:"password"`
	SigninMethod string `json:"signinMethod"`
}
