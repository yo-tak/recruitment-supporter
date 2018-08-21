package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
)

func main() {
	// extract user information
	recruitPageURL, companyName, userid, password, signinMethod, err := extractArgs()
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
	if err := page.Navigate(recruitPageURL); err != nil {
		log.Fatalf("Failed to navigate to main page: %v", err)
	}
	if err := authenticate(page, signinMethod, userid, password); err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	if err := page.Navigate(recruitPageURL + "/companies/" + companyName + "/projects"); err != nil {
		log.Fatalf("Failed to open company page: %v", err)
	}
	// TODO I just couldn't afford the time to handle error here
	supportOffers(page)
	log.Println("Finished! Thank you so much for your support!!")
}

func extractArgs() (string, string, string, string, string, error) {
	args := os.Args
	if len(args) < 5 {
		return "", "", "", "", "", errors.New("not enough args; recruit page url, company name, userId, and password are required")
	}
	recruitPageURL := args[1]
	companyName := args[2]
	userid := args[3]
	password := args[4]
	signinMethod := ""
	if len(args) > 5 {
		signinMethod = args[5]
	}
	return recruitPageURL, companyName, userid, password, signinMethod, nil
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

	time.Sleep(2 * time.Second) // wait until twitter window is ready
	// close twitter window
	// TODO should improve logic to handle windows
	page.NextWindow()                    // move to new Twitter window
	page.CloseWindow()                   // close Twitter window
	page.Session().SetWindow(mainWindow) // move back to recruitment page window

	time.Sleep(1 * time.Second) // wait until wantedly window is ready again

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
