package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func init() {
	go downloadDoubleThread()
}

//DoubleIndexFilter puts the domains inside file
func DoubleIndexFilter(url string) error {

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

	struc := regexp.MustCompile("^([0-9]+[.].+)[ ]+(.+)$")

	for _, a := range dlines {

		if struc.MatchString(a) {
			k := struc.FindStringSubmatch(a)
			DomainKill(k[2])
		}
	}

	fmt.Println("Ingested url: ", url)

	return err

}

func getDoubleFilters() {

	s := strings.Split(ZabovDoubleBL, ",")

	for _, a := range s {
		DoubleIndexFilter(a)
	}

}

func downloadDoubleThread() {
	fmt.Println("Starting updater of DOUBLE lists")
	for {
		getDoubleFilters()
		time.Sleep(12 * time.Hour)
	}

}
