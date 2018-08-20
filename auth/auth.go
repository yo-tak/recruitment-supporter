package auth

import (
	"log"

	"github.com/sclevine/agouti"
)

func Authenticate(page *agouti.Page, signinMethod string, username string, password string) error {

	switch signinMethod {
	case "facebook":
		if err := page.Find("ul.nav .ui-show-modal").Click(); err != nil {
			log.Printf("Failed to open login dialog: %v\n", err)
			return err
		}
		page.Find(".login-button.facebook").Click()
		if err := page.FindByID("email").Fill(username); err != nil {
			log.Printf("Failed to fill in user email: %v\n", err)
		}
		if err := page.FindByID("pass").Fill(password); err != nil {
			log.Printf("Failed to fill in user password: %v\n", err)
		}
		if err := page.FindByID("loginbutton").Click(); err != nil {
			log.Printf("Failed to click login button: %v\n", err)
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
		return nil
	}
}
