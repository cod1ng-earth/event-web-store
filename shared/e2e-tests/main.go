// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

// https://github.com/chromedp/chromedp/issues/82#issuecomment-312022893

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/cdproto/emulation"
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

	// add a timeout
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()
	width, height := 1024, 768
	chromedp.Run(ctx, emulation.SetDeviceMetricsOverride(int64(width), int64(height), 1.0, false))
	log.Printf("staring and resizing chrome took %v", time.Since(start))

	start = time.Now()
	picture := &[]byte{}
	var h2 string
	chromedp.Run(
		ctx,
		//5*time.Second,
		chromedp.Navigate(`http://localhost:8000/`),
		//		chromedp.Reload(),
		chromedp.Click(`#main > div > ul > li:nth-child(2) > span:nth-child(3) > button`, chromedp.NodeVisible), // add product to cart
		chromedp.Click(`body > div > header > div > button:nth-child(4)`),                                       // go to cart
		chromedp.Click(`#main > div > ul > li:nth-child(2) > button`),                                           // order now
		chromedp.Text(`#main > h2`, &h2),
		chromedp.CaptureScreenshot(picture),
		screenshotSave("cart_ordered.png", picture))
	log.Printf("test took %v", time.Since(start))
	log.Printf("h2: %v", h2)

	start = time.Now()
	picture = &[]byte{}
	var price string
	chromedp.Run(
		ctx,
		//5*time.Second,
		chromedp.Navigate(`http://localhost:8000/`),
		//		chromedp.Reload(),
		chromedp.Click(`#main > div > ul > li:nth-child(1) > span:nth-child(2) > button`, chromedp.NodeVisible), // go to detail page
		chromedp.Text(`#main > div > div.mdl-cell.mdl-cell--4-col > span.mdl-typography--display-1.custom-detail-block`, &price),
		chromedp.CaptureScreenshot(picture),
		screenshotSave("detail_page.png", picture))
	log.Printf("test took %v", time.Since(start))
	log.Printf("price: %v", price)
}

/*
func runActions(ctx context.Context d time.Duration, actions ...chromedp.Action) error {
	//ctx, cancel = context.WithTimeout(ctx, d)
	//defer cancel()

	width, height := 1024, 768
	actions = prepend(emulation.SetDeviceMetricsOverride(int64(width), int64(height), 1.0, false), actions)

	return chromedp.Run(ctx, actions...)
}

func prepend(e chromedp.Action, ls []chromedp.Action) []chromedp.Action {
	return append([]chromedp.Action{e}, ls...)
}
*/

func screenshotSave(fileName string, buf *[]byte) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		//return nil
		log.Printf("Write %v", fileName)
		return ioutil.WriteFile(fileName, *buf, 0644)
	}
}
