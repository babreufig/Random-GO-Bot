package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func createNewWTRoom() (string, error) {

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	//body := strings.NewReader(`authenticity_token=MTYni%2BjqL4ubvwM5UABXsCmsuK3en9TanUrG69ns9eXjs%2BQ2e80VZTkRzg2AorN3Nsv76Mz2Fx%2FDBKH4EufjoQ%3D%3D`)
	body := strings.NewReader("")
	req, err := http.NewRequest("POST", "https://www.watch2gether.com/rooms/create", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authority", "www.watch2gether.com")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Ch-Ua", "Google Chrome 75")
	req.Header.Set("Origin", "https://www.watch2gether.com")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Origin-Policy", "0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Set("Referer", "https://www.watch2gether.com/")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br") To make it easier to decode don't use compression
	req.Header.Set("Accept-Language", "de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7")
	//req.Header.Set("Cookie", "w2g_session_id=80b4ed334e453584ad1250825fecbff5; w2g_cookies_accepted=1; paddlejs_checkout_variant={\"inTest\":true,\"controlGroup\":false,\"isForced\":false,\"variant\":\"multipage-radio-payment\"}; paddlejs_campaign_referrer=www.watch2gether.com; w2glang=de; _pk_ref.1.b9aa=%5B%22%22%2C%22%22%2C1562077005%2C%22https%3A%2F%2Fwww.google.com%2F%22%5D; _pk_ses.1.b9aa=1; _pk_id.1.b9aa=58238e481d09ea17.1560021959.12.1562077112.1562077005.")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyText, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyText)
	roomID := bodyString[70:88]

	if `<html><body>You are being <a href="https://www.watch2gether.com/rooms/` != bodyString[:70] || len(bodyString) != 127 {
		return "", errors.New("unknown response by Watch2Gether api")
	}

	return roomID, nil
}
