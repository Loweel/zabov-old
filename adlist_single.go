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
	fmt.Println("Killfile size: ", len(myBody))

	dlines := strings.Split(myBody, "\n")

	incrementStats("PlainLines " + durl, int64(len(dlines)))
	
	for _, a := range dlines {

		b := strings.Fields(a)

		if len(b) < 1 {
			continue
		}

		ur, _ := url.Parse("http://" + b[0])

		if ur.IsAbs() {
			DomainKill(ur.Hostname(), durl)
		} else {
			incrementStats("Malformed Single",1 )
			
			
		}

	}

	incrementStats("SourcesList",1 )
	

	return err

}

func getSingleFilters() {

	s := strings.Split(ZabovSingleBL, ",")

	for _, a := range s {
		SingleIndexFilter(a)
	}

}

func downloadThread() {
	fmt.Println("Starting updater of SINGLE lists, each (hours): ", ZabovCacheTTL)
	for {
		getSingleFilters()
		time.Sleep(time.Duration(ZabovCacheTTL) * time.Hour)
	}

}
