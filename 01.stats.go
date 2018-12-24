package main

import (
	"time"
	"sync"
	"encoding/json"
	"fmt"
				
)

//ZabovStats is used to keep statistics to print
var ZabovStats map[string]int64

//StatMutex is for avoid race condition on the map
var StatMutex = new(sync.Mutex)

func init(){


	ZabovStats = make(map[string]int64)

	fmt.Println("Initializing stats engine.")
	go reportPrintThread()


}


func statsPrint(){
	StatMutex.Lock()
	fmt.Println()
	stat,_ := json.Marshal(ZabovStats)
	fmt.Println(string(stat))
	StatMutex.Unlock()
	fmt.Println()

}


func incrementStats(key string, value int64){

	StatMutex.Lock()
	ZabovStats[key] += value 
	StatMutex.Unlock()


}

func setstatsvalue(key string, value int64){

	StatMutex.Lock()
	ZabovStats[key] = value
	StatMutex.Unlock()


}




func reportPrintThread(){
	for {
	statsPrint()
	time.Sleep(2 * time.Minute)
	}
}

