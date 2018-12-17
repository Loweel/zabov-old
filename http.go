package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func init() {
	downloadThread()
}

//SingleIndexFilter puts the domains inside file
func SingleIndexFilter(url string) error {

	fmt.Println("Retrieving killfile from: ", url)

	var myBody string
	var bodyBytes []byte
	var err error

	// Get the data
	resp, err := http.Get(url)
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

	fmt.Println("Number of lines: ", len(dlines))

	for _, a := range dlines {

		if strings.Contains(a, "#") == false {
			DomainKill(a)
		}
	}

	fmt.Println("Ingested url: ", url)

	return err

}

func getSingleFilters() {

	s := strings.Split(ZabovSingleBL, ",")

	for _, a := range s {
		SingleIndexFilter(a)
	}

}

func downloadThread() {
	for {
		getSingleFilters()
		time.Sleep(12 * time.Hour)
	}

}
