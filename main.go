package main

import (
	"errors"
	"log"
	"os"

	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
)

func main() {
	// extract user information
	recruitPageUrl, companyName, userid, password, err := extractArgs()
	if err != nil {
		log.Fatalf("Failed to extract necessary args: %v", err)
	}

	driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver: %v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Fatalf("Failed to open page: %v", err)
	}

	// open wantedly main page
	if err := page.Navigate(recruitPageUrl); err != nil {
		log.Fatalf("Failed to navigate to main page: %v", err)
	}
	if err := login(page, userid, password); err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	if err := page.Navigate(recruitPageUrl + "/companies/" + companyName + "/projects"); err != nil {
		log.Fatalf("Failed to open company page: %v", err)
	}
	// TODO I just couldn't afford the time to handle error here
	supportOffers(page)
	log.Println("Finished! Thank you so much for your support!!")
}

func extractArgs() (string, string, string, string, error) {
	args := os.Args
	if len(args) < 5 {
		return "", "", "", "", errors.New("not enough args; I need recruit page url, company name, userId, and password")
	}
	recruitPageUrl := args[1]
	companyName := args[2]
	userid := args[3]
	password := args[4]
	return recruitPageUrl, companyName, userid, password, nil
}

func login(page *agouti.Page, name string, password string) error {
	// TODO is there any smarter way to handle the errors?
	if err := page.Find("ul.nav .ui-show-modal").Click(); err != nil {
		log.Printf("Failed to open login dialog: %v\n", err)
		return err
	}
	if err := page.FindByID("login_user_email").Fill(name); err != nil {
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

func supportOffers(page *agouti.Page) {
	supportIcons := page.AllByClass("wt-icon-support")
	elems, err := supportIcons.Elements()
	if err != nil {
		log.Fatalf("Failed to extract support icons: %v", err)
	}
	mainWindow, _ := page.Session().GetWindow()
	for _, elem := range elems {
		supportOffer(elem, page, mainWindow)
	}
}

func supportOffer(elem *api.Element, page *agouti.Page, mainWindow *api.Window) {
	// open support dialog and click support by twitter
	elem.Click()
	supportDialog := page.FindByID("wtd-modal-portal__default")
	if err := supportDialog.FindByButton("Twitterで応援").Click(); err != nil {
		log.Fatalf("Failed to click support button for elem: %v, err: %v", elem, err)
	}

	// close twitter window
	// TODO should improve logic to handle windows
	// move to new Twitter window
	page.NextWindow()
	// close Twitter window
	page.CloseWindow()
	// move back to recruitment page window
	page.Session().SetWindow(mainWindow)

	// close support dialog. when it's your first time to supported the offer,
	// the caption will be "応援しない". otherwise it will be "閉じる""
	closeButton := supportDialog.FindByButton("応援しない")
	if count, err := closeButton.Count(); err != nil || count == 0 {
		closeButton = supportDialog.FindByButton("閉じる")
	}
	if err := closeButton.Click(); err != nil {
		log.Fatalf("Failed to click close button for elem: %v, err: %v", elem, err)
	}
}
