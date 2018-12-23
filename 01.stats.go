package main

import (
	"time"
	
	
	"fmt"
				
)

//ZabovStats is used to keep statistics to print
var ZabovStats map[string]int64


func init(){


	ZabovStats = make(map[string]int64)

	fmt.Println("Initializing stats engine.")
	go reportPrintThread()


}


func statsPrint(){
	fmt.Println()
	fmt.Println("Usage Stats: ")
	for key,value := range ZabovStats {
	
		fmt.Printf("%s : %d\n", key, value)

	}
	fmt.Println()

}

func reportPrintThread(){
	statsPrint()
	time.Sleep(2 * time.Minute)

}

