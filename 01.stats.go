package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

//ZabovStats is used to keep statistics to print
var ZabovStats map[string]int64

var increment *list.List
var overwrite *list.List

//StatMutex is for avoid race condition on the map
var StatMutex = new(sync.Mutex)

func init() {

	increment = list.New()
	overwrite = list.New()

	ZabovStats = make(map[string]int64)

	fmt.Println("Initializing stats engine.")
	go reportPrintThread()
	go incrementThread()
	go overwriteThread()

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

	s := fmt.Sprintf("%s|%d", key, value)

	increment.PushBack(s)

}

func setstatsvalue(key string, value int64) {

	s := fmt.Sprintf("%s|%d", key, value)

	overwrite.PushBack(s)

}

func reportPrintThread() {
	for {
		statsPrint()
		time.Sleep(2 * time.Minute)
	}
}

func incrementThread() {

	fmt.Println("Starting Stats Increment Thread")

	for {

		for increment.Len() > 0 {
			e := increment.Front() // First element

			a := strings.Split(fmt.Sprint(e.Value), "|")
			inc, _ := strconv.ParseInt(a[1], 10, 64)
			StatMutex.Lock()
			ZabovStats[a[0]] += inc
			StatMutex.Unlock()
			increment.Remove(e) // Dequeue
		}
		time.Sleep(10 * time.Millisecond)
	}

}

func overwriteThread() {

	fmt.Println("Starting Stats Overwrite Thread")

	for {

		for overwrite.Len() > 0 {
			e := overwrite.Front() // First element

			a := strings.Split(fmt.Sprint(e.Value), "|")
			inc, _ := strconv.ParseInt(a[1], 10, 64)
			StatMutex.Lock()
			ZabovStats[a[0]] = inc
			StatMutex.Unlock()
			overwrite.Remove(e) // Dequeue
		}
		time.Sleep(10 * time.Millisecond)
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
