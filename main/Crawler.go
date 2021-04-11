package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

func GetHTTPClient() *http.Client {
	cl := &http.Client{}
	return cl
}

const (
	addToCart = "add-to-cart-button"
	buyNow    = "buy-now-button"
)

func sendMail(subject string, content string) {

	from := "icheckforps5allday@gmail.com"
	pass := "nEXa5iRBxiCh3UY"

	to := "manuel.welte94@gmail.com"

	//msg := "Subject: Control (No news)\n" + "Test email"
	msg := "Subject: " + subject + "\n" + content

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		panic(err)

	}
	fmt.Println("Sent email")
}

func main() {

	url := "https://www.amazon.de/-/en/9843344/dp/B01733YMKU/ref=sr_1_3?dchild=1&keywords=bloodborne&qid=1618157407&s=videogames&sr=1-3"
	url = "https://www.amazon.de/-/en/dp/B08H93ZRK9/ref=sr_1_2?crid=3U3AMAJVF68AS&dchild=1&keywords=playstation+5&qid=1618162759&s=videogames&sprefix=play%2Cvideogames%2C198&sr=1-2"
	cl := GetHTTPClient()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0")

	resp, err := cl.Do(req)
	if resp == nil {
		panic(err)
	}
	bodyOld, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(bodyOld))
	if bodyOld == nil {
		panic(err)
	}

	sendMail("Control", "Starting to poll" )

	for true {
		resp, err = cl.Do(req)

		if resp == nil {
			// Don't quit (sleep & restart)
			fmt.Println(err)
			time.Sleep(1 * time.Minute)
			continue
		}

		bodyNew, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// Restart?
			fmt.Println(err)
			time.Sleep(1 * time.Minute)
			continue
		}

		if !bytes.Equal(bodyNew, bodyOld) {
			fmt.Println("Content changed")
			handleDiff(&bodyNew)
		} else {
			fmt.Println("Same old")
		}

		bodyOld = bodyNew

		time.Sleep(5 * time.Second)
	}
}

func handleDiff(contentB *[]byte) {
	content := string(*contentB)
	if strings.Contains(content, addToCart) {
		sendMail("Still active", "am still here")
		fmt.Println("Contains button with id:", addToCart)
	}
	if strings.Contains(content, buyNow) {
		sendMail("ALERT: AVAILABLE", buyNow + " seems to be active")
		fmt.Println("Contains button with id:", buyNow)
	}
}
