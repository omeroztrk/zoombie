package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/playwright-community/playwright-go"
)

func main() {
	hless_p := flag.Bool("hl", true, "Headless")
	name_p := flag.String("n", "", "Name of the bot participant")
	zoom_link_p := flag.String("zl", "", "Zoom invite link")
	flag.Parse()

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
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	page.Context().GrantPermissions([]string{`microphone`, `camera`})
	if _, err = page.Goto(meetLink); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	cookieBox, _ := page.WaitForSelector(`[aria-label="Privacy"]`)
	fmt.Println("Cookie thing popped up")
	if closeCookie, err := cookieBox.QuerySelector(`[aria-label="Close"]`); err == nil {
		closeCookie.Click()
	}

	nameFiled, _ := page.QuerySelector(`[placeholder="Your Name"]`)
	nameFiled.Fill(name)

	joinBtn, _ := page.QuerySelector(`#joinBtn`)
	joinBtn.Click()

	page.WaitForSelector(`[aria-label="open the chat pane"]`)
	page.Click(`[aria-label="open the chat pane"]`)

	if mute, err := page.QuerySelector(`button:has-text("Mute")`); err == nil && mute != nil {
		mute.Click()
	}

	if svideo, err := page.QuerySelector(`button:has-text("Stop Video")`); err == nil && svideo != nil {
		svideo.Click()
	}

	fmt.Println("Joined")
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
