package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func init() {
	go downloadThread()
}

//SingleIndexFilter puts the domains inside file
func SingleIndexFilter(durl string) error {

	if _, urlErr := url.ParseRequestURI(durl); urlErr != nil {
		return urlErr
	}

	fmt.Println("Retrieving killfile from: ", durl)

	var myBody string
	var bodyBytes []byte
	var err error

	// Get the data
	resp, err := http.Get(durl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 { // OK
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Problem downloading: ", err.Error())
		}

	} else {
		bodyBytes = []byte{}
	}

	myBody = string(bodyBytes)
	fmt.Printf("%s Killfile size: %d\n", durl, len(myBody))

	dlines := strings.Split(myBody, "\n")

	go setstatsvalue("PlainLines "+durl, int64(len(dlines)))

	for _, a := range dlines {

		b := strings.Fields(a)

		if len(b) != 1 {
			continue
		}

		ur, urStrErr := url.ParseRequestURI("http://" + b[0])
		if urStrErr != nil {
			go incrementStats("Malformed Domains", 1)
			return urStrErr
		}

		if ur.IsAbs() {
			DomainKill(ur.Hostname(), durl)
		} else {
			fmt.Print(b)
			go incrementStats("Malformed Domains", 1)

		}

	}

	go incrementStats("SourcesListUrl", 1)

	return err

}

func getSingleFilters() {

	s := fileByLines(ZabovSingleBL)

	for _, a := range s {
		go SingleIndexFilter(a)
	}

}

func downloadThread() {
	fmt.Println("Starting updater of SINGLE lists, each (hours): ", ZabovKillTTL)
	for {
		getSingleFilters()
		time.Sleep(time.Duration(ZabovKillTTL) * time.Hour)
	}

}
