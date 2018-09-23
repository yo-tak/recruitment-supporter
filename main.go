package main

import (
	"log"
	"os"
	"time"

	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
)

type statusCode int

const (
	success          statusCode = iota
	unsufficientCfg  statusCode = iota
	driverError      statusCode = iota
	signinFailed     statusCode = iota
	recruitPageError statusCode = iota
	otherFailure     statusCode = iota
)

func main() {
	result := execSupport()
	os.Exit(int(result))
}

func execSupport() statusCode {

	// extract user information
	cfg, err := getCfg()
	if err != nil {
		log.Printf("Failed to extract necessary args: %v", err)
		return unsufficientCfg
	}

	driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		log.Printf("Failed to start driver: %v", err)
		return driverError
	}
	defer driver.Stop()

	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Printf("Failed to open page: %v", err)
		return driverError
	}
	defer page.CloseWindow()

	// open wantedly main page
	if err := page.Navigate(cfg.URL); err != nil {
		log.Printf("Failed to navigate to main page: %v", err)
		return recruitPageError
	}
	if err := authenticate(page, cfg.SigninMethod, cfg.Mail, cfg.Password); err != nil {
		log.Printf("Failed to sign in: %v", err)
		return signinFailed
	}

	if err := page.Navigate(cfg.URL + "/companies/" + cfg.Company + "/projects"); err != nil {
		log.Printf("Failed to open company page: %v", err)
		return recruitPageError
	}
	// TODO I just couldn't afford the time to handle error here
	supportOffers(page)
	log.Println("Finished! Thank you so much for your support!!")
	return success
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

	time.Sleep(1500 * time.Millisecond) // wait until twitter window is ready
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
