// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		func(a *chromedp.ExecAllocator) {
			chromedp.Flag("headless", false)(a)
			// Like in Puppeteer.
			chromedp.Flag("hide-scrollbars", false)(a)
			chromedp.Flag("mute-audio", false)(a)
		})

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	var example string
	start := time.Now()
	err := chromedp.Run(ctx,
		chromedp.Navigate(`http://localhost:8000/`),
		// wait for footer element is visible (ie, page is loaded)
		//		chromedp.WaitVisible(`#main > div > ul > li:nth-child(1) > span:nth-child(3) > button`),
		// find and click "Expand All" link
		chromedp.Click(`#main > div > ul > li:nth-child(2) > span:nth-child(3) > button`, chromedp.NodeVisible),
		chromedp.Click(`body > div > header > div > button:nth-child(4)`),
		//		chromedp.Click(`#main > ul > li.mdl-list__item.mdl-list__item--two-line > span:nth-child(2) > button:nth-child(3)`, chromedp.NodeVisible),
		//		chromedp.Click(`#main > ul > li.mdl-list__item.mdl-list__item--two-line > span:nth-child(2) > button:nth-child(3)`),
		chromedp.Click(`#main > ul > li:nth-child(2) > button`),

		//		chromedp.WaitVisible(`#main > h1`),
		//		chromedp.Value(`#main > h1`, &example),
		// retrieve the value of the textarea
		chromedp.Text(`#main > h1`, &example),
	)
	log.Printf("test took %v", time.Since(start))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Go's time.After example:\n%s", example)
}
