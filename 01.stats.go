package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type send struct {
	Payload string
	Number  int64
}

//ZabovStats is used to keep statistics to print
var ZabovStats map[string]int64

var increment chan send
var overwrite chan send

//StatMutex is for avoid race condition on the map
var StatMutex = new(sync.Mutex)

func init() {

	increment = make(chan send, 1024)
	overwrite = make(chan send, 1024)

	ZabovStats = make(map[string]int64)

	fmt.Println("Initializing stats engine.")
	go reportPrintThread()
	go statsThread()
}

func statsPrint() {
	StatMutex.Lock()
	fmt.Println()
	stat, _ := json.Marshal(ZabovStats)
	fmt.Println(jsonPrettyPrint(string(stat)))
	StatMutex.Unlock()
	fmt.Println()

}

func incrementStats(key string, value int64) {

	var s send

	s.Payload = key
	s.Number = value

	increment <- s

}

func setstatsvalue(key string, value int64) {

	var s send

	s.Payload = key
	s.Number = value

	overwrite <- s

}

func reportPrintThread() {
	for {
		statsPrint()
		time.Sleep(2 * time.Minute)
	}
}

func statsThread() {

	fmt.Println("Starting Statistical Collection Thread")

	for {
		select {
		case ovrw := <-overwrite:
			StatMutex.Lock()
			ZabovStats[ovrw.Payload] = ovrw.Number
			StatMutex.Unlock()
		case incr := <-increment:
			StatMutex.Lock()
			ZabovStats[incr.Payload] += incr.Number
			StatMutex.Unlock()
		}

	}
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}
