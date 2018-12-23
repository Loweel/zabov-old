package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func init() {
	go downloadDoubleThread()
}

//DoubleIndexFilter puts the domains inside file
func DoubleIndexFilter(durl string) error {

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

	ZabovStats["HostsLines " + durl] = int64(len(dlines))

	for _, a := range dlines {

		k := strings.Fields(a)

		if len(k) > 1 {
			if net.ParseIP(k[0]) != nil {
				DomainKill(strings.Trim(k[1], " "), durl)
			}
		} else {
			ZabovStats["Malformed Host" + durl]++
			
		}

	}

	ZabovStats["SourcesHost"]++

	return err

}

func getDoubleFilters() {

	s := strings.Split(ZabovDoubleBL, ",")

	for _, a := range s {
		DoubleIndexFilter(a)
	}

}

func downloadDoubleThread() {
	fmt.Println("Starting updater of DOUBLE lists, each (hours):", ZabovCacheTTL)
	for {
		getDoubleFilters()
		time.Sleep(time.Duration(ZabovCacheTTL) * time.Hour)
	}

}
