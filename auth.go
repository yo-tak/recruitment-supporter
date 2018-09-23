package main

import (
	"errors"
	"log"
	"strings"

	"github.com/sclevine/agouti"
)

func authenticate(page *agouti.Page, signinMethod string, username string, password string) error {

	switch signinMethod {
	case "facebook":
		if err := page.Find("ul.nav .ui-show-modal").Click(); err != nil {
			log.Printf("Failed to open login dialog: %v\n", err)
			return err
		}
		page.Find(".login-button.facebook").Click()
		if err := page.FindByID("email").Fill(username); err != nil {
			log.Printf("Failed to fill in user email: %v\n", err)
			return err
		}
		if err := page.FindByID("pass").Fill(password); err != nil {
			log.Printf("Failed to fill in user password: %v\n", err)
			return err
		}
		if err := page.FindByID("loginbutton").Click(); err != nil {
			log.Printf("Failed to click login button: %v\n", err)
			return err
		}
		url, err := page.URL()
		if err != nil {
			log.Printf("Failed to get URL after signin attempt: %v\n", err)
			return err
		}
		if strings.Contains(url, "facebook") {
			return errors.New("Failed to signin; mail address and/or password is wrong")
		}
		return nil
	default: // default case supports signin that does not use any third party
		// TODO is there any smarter way to handle the errors?
		if err := page.Find("ul.nav .ui-show-modal").Click(); err != nil {
			log.Printf("Failed to open login dialog: %v\n", err)
			return err
		}
		if err := page.FindByID("login_user_email").Fill(username); err != nil {
			log.Printf("Failed to fill in user email: %v\n", err)
			return err
		}
		if err := page.FindByID("login_user_password").Fill(password); err != nil {
			log.Printf("Failed to fill in password: %v\n", err)
			return err
		}
		if err := page.Find("#login_new_user input[type=\"submit\"]").Click(); err != nil {
			log.Printf("Failed to click login page: %v\n", err)
			return err
		}
		url, err := page.URL()
		if err != nil {
			log.Printf("Failed to get URL after signin attempt: %v\n", err)
			return err
		}
		if strings.Contains(url, "sign_in") {
			return errors.New("mail address and/or password is wrong")
		}
		return nil
	}
}
