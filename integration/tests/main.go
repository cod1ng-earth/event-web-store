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

	start := time.Now()
	str := runChrome()
	log.Printf("test took %v", time.Since(start))
	log.Printf("Go's time.After example:\n%s", str)

}

func runChrome() string {
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
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	picture := &[]byte{}
	width, height := 1024, 768
	var example string

	err := chromedp.Run(ctx,
		emulation.SetDeviceMetricsOverride(int64(width), int64(height), 1.0, false),
		chromedp.Navigate(`http://localhost:8000/`),
		chromedp.Click(`#main > div > ul > li:nth-child(2) > span:nth-child(3) > button`, chromedp.NodeVisible), // add product to cart
		chromedp.Click(`body > div > header > div > button:nth-child(4)`),                                       // go to cart
		chromedp.Click(`#main > div > ul > li:nth-child(2) > button`),                                           // order now
		chromedp.Text(`#main > h2`, &example),
		chromedp.CaptureScreenshot(picture),
		screenshotSave("abc.png", picture),
	)
	if err != nil {
		log.Printf("%v", err)
	}

	return example
}

func screenshotSave(fileName string, buf *[]byte) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		log.Printf("Write %v", fileName)
		return ioutil.WriteFile(fileName, *buf, 0644)
	}
}
