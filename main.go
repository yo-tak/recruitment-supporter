package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
)

type statusCode int

const (
	success statusCode = iota
	unsufficientCfg
	driverError
	signinFailed
	recruitPageError
	otherFailure
)

type lang string

const (
	english  lang = "en"
	japanese lang = "ja"
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

func getLocale(page *agouti.Page) (lang, error) {
	cookies, err := page.GetCookies()
	if err != nil {
		return "", err
	}
	for _, cookie := range cookies {
		if cookie.Name == "locale" {
			switch cookie.Value {
			case "en":
				return english, nil
			case "ja":
				return japanese, nil
			default:
				return "", errors.New(fmt.Sprintf("no language found for locale: %v", cookie.Value))
			}
		}
	}
	log.Println("no locale set to cookie. set language to japanese.")
	return japanese, nil
}

func supportOffer(elem *api.Element, page *agouti.Page, mainWindow *api.Window) error {
	language, err := getLocale(page)
	if err != nil {
		return err
	}
	// open support dialog and click support by twitter
	elem.Click()
	supportDialog := page.FindByID("wtd-modal-portal__default")

	supportButtonCaption := ""
	if language == japanese {
		supportButtonCaption = "Twitterで応援"
	} else {
		supportButtonCaption = "Recommend on Twitter"
	}
	supportButton := supportDialog.FindByButton(supportButtonCaption)
	if err := supportButton.Click(); err != nil {
		log.Printf("Failed to click support button for elem: %v, err: %v", elem, err)
		return err
	}
	proceedTwitterWindow(page)

	// TODO should improve logic to handle windows
	if err := page.Session().SetWindow(mainWindow); err != nil { // move back to recruitment page window
		return err
	}
	time.Sleep(1 * time.Second) // wait until wantedly window is ready again

	// close support dialog. when it's your first time to supported the offer,
	// the caption will be "応援しない". otherwise it will be "閉じる""
	closeButtonCaption := ""
	if language == japanese {
		closeButtonCaption = "閉じる"
	} else {
		closeButtonCaption = "Close"
	}
	closeButton := supportDialog.FindByButton("応援しない")
	if count, err := closeButton.Count(); err != nil || count == 0 {
		closeButton = supportDialog.FindByButton(closeButtonCaption)
	}
	if err := closeButton.Click(); err != nil {
		log.Printf("Failed to click close button for elem: %v, err: %v", elem, err)
		return err
	}
	return nil
}

func proceedTwitterWindow(page *agouti.Page) error {
	if err := page.NextWindow(); err != nil { // move to new Twitter window
		return err
	}
	twitterReady := make(chan bool)
	go func() {
		for trial := 0; trial < 20; trial++ {
			time.Sleep(1000 * time.Millisecond)
			tweetForm := page.Find("input[type='submit']")
			if tweetForm != nil {
				twitterReady <- true
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
		twitterReady <- false
	}()
	if !<-twitterReady {
		return errors.New("failed to open twitter form")
	}
	if err := page.CloseWindow(); err != nil { // close Twitter window
		return err
	}
	return nil
}
