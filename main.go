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
	// deligate procedure to execSupport func so that deferred process will be executed
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

	// open main page
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
	if err := supportOffers(page); err != nil {
		log.Printf("Failed in supporting process: %v", err)
		return recruitPageError
	}
	log.Println("Finished! Thank you so much for your support!!")
	return success
}

func supportOffers(page *agouti.Page) error {
	supportIcons := page.AllByClass("wt-icon-support")
	elems, err := supportIcons.Elements()
	if err != nil {
		log.Printf("Failed to extract support icons: %v\n", err)
		return err
	}
	mainWindow, err := page.Session().GetWindow()
	if err != nil {
		log.Printf("Failed to get new window: %v\n", err)
		return err
	}
	for _, elem := range elems {
		if err := supportOffer(elem, page, mainWindow); err != nil {
			return err
		}
	}
	return nil
}

func supportOffer(elem *api.Element, page *agouti.Page, mainWindow *api.Window) error {
	// open support dialog and click support by twitter
	elem.Click()
	supportDialog := page.FindByID("wtd-modal-portal__default")
	if err := supportDialog.FindByButton("Twitterで応援").Click(); err != nil {
		log.Printf("Failed to click support button for elem: %v, err: %v", elem, err)
		return err
	}

	time.Sleep(1500 * time.Millisecond) // wait until twitter window is ready
	// close twitter window
	// TODO should improve logic to handle windows
	if err := page.NextWindow(); err != nil { // move to new Twitter window
		return err
	}
	if err := page.CloseWindow(); err != nil { // close Twitter window
		return err
	}
	if err := page.Session().SetWindow(mainWindow); err != nil { // move back to recruitment page window
		return err
	}
	time.Sleep(1 * time.Second) // wait until wantedly window is ready again

	// close support dialog. when it's your first time to supported the offer,
	// the caption will be "応援しない". otherwise it will be "閉じる""
	closeButton := supportDialog.FindByButton("応援しない")
	if count, err := closeButton.Count(); err != nil || count == 0 {
		closeButton = supportDialog.FindByButton("閉じる")
	}
	if err := closeButton.Click(); err != nil {
		log.Printf("Failed to click close button for elem: %v, err: %v", elem, err)
		return err
	}
	return nil
}
