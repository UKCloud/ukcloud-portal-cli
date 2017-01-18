package command

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type CookiesCollection struct {
	Collection []*http.Cookie
}

var cookieCollection = new(CookiesCollection)

func auth(email string, password string) int {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	if len(email) < 1 || len(password) < 1 {
		return 1
	}

	myClient := http.Client{Jar: jar, Timeout: 100 * time.Second}

	var jsonStr = []byte(`{"email": "` + email + ` ", "password": "` + password + `"}`)

	url := "https://portal.skyscapecloud.com/api/authenticate.json"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	r, err := myClient.Do(req)

	if r.StatusCode != 201 {
		return 1
	}

	cookieCollection.Collection = r.Cookies()
	return 0
}

func getJson(myUrl string, target interface{}) error {

	jar, _ := cookiejar.New(nil)

	u, err := url.Parse(myUrl)
	jar.SetCookies(u, cookieCollection.Collection)

	tr := &http.Transport{
		DisableCompression: true,
	}

	myClient := http.Client{Jar: jar, Timeout: 100 * time.Second, Transport: tr}

	r, err := myClient.Get(myUrl)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err := json.Unmarshal(body, target); err != nil {
		log.Println(err)
	}

	return json.Unmarshal(body, target)
}
