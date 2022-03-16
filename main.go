package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

func main() {
	hless_p := flag.Bool("hl", true, "Headless")
	name_p := flag.String("n", "", "Name of the bot participant")
	zoom_link_p := flag.String("zl", "", "Zoom invite link")
	verbosity_p := flag.Int("v", 1, "Verbosity 1-3")

	flag.Parse()

	if *verbosity_p > 1 {
		log.Println("Verbosity: ", *verbosity_p)
	}

	invite := *zoom_link_p
	name := *name_p

	if name == "" || invite == "" {
		log.Fatalln("You shall not pass. Enter all params")
	}

	meetLink := InviteToBrowserLink(invite)

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: hless_p})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	if *verbosity_p > 1 {
		log.Println("Chromium launched")
	}

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	page.Context().GrantPermissions([]string{`microphone`, `camera`})
	if *verbosity_p > 1 {
		log.Println("Permissions granted")
	}

	if _, err = page.Goto(meetLink); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	if *verbosity_p > 1 {
		log.Println("Went to: ", meetLink)
	}

	cookieBox, _ := page.WaitForSelector(`[aria-label="Privacy"]`)
	if *verbosity_p > 1 {
		log.Println("Cookie thing popped up")
	}

	if closeCookie, err := cookieBox.QuerySelector(`[aria-label="Close"]`); err == nil {
		closeCookie.Click()
	}

	if *verbosity_p > 1 {
		log.Println("Cookie thing closed")
	}

	nameFiled, _ := page.QuerySelector(`[placeholder="Your Name"]`)
	nameFiled.Fill(name)
	if *verbosity_p > 1 {
		log.Println("Name entered")
	}

	joinBtn, _ := page.QuerySelector(`#joinBtn`)
	joinBtn.Click()
	if *verbosity_p > 1 {
		log.Println("Join button clicked")
	}

	if agreeButton, err := page.QuerySelector(`#wc_agree1`); err == nil && agreeButton != nil {
		agreeButton.Click()
		log.Println("Agree1 clicked")
	}

	if agreeButton, err := page.QuerySelector(`#wc_agree2`); err == nil && agreeButton != nil {
		agreeButton.Click()
		log.Println("Agree2 clicked")
	}

	page.WaitForSelector(`[aria-label="open the chat pane"]`)
	page.Click(`[aria-label="open the chat pane"]`)
	log.Println("Joined")

	// delayed screenshot
	if *verbosity_p > 2 {
		go func() { // since that, rest can continue
			time.Sleep(20 * time.Second) //lets screenshot after 20 secs
			log.Println("Screenshotting after login")
			path := "Delayed After Login.png"
			fullPage := true

			page.Screenshot(playwright.PageScreenshotOptions{
				Path:     &path,
				FullPage: &fullPage,
			})
		}()
	}
	log.Println("Continuing")

	for {
		var line string
		reader := bufio.NewReader(os.Stdin)
		line, _ = reader.ReadString('\n')
		if line == "exit\n" { // be careful. a small typo may lead catastrophic misunderstandings
			break
		}
		page.Click(`textarea[type="text"]`)
		page.Fill(`textarea[type="text"]`, line)
		page.Press(`textarea[type="text"]`, `Enter`)
		fmt.Println("Sent: ", line)
	}

	browser.Close()
	pw.Stop()
}

func InviteToBrowserLink(invite string) string {
	return strings.Replace(invite, `/j/`, `/wc/join/`, -1)
}
