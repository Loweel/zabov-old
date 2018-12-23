package main

import (
	
	"os"
	"bufio"
	"fmt"
				
)

func init(){

fmt.Println("Ingesting local hosts file")
ingestLocalBlacklist()
	
	
}

func ingestLocalBlacklist(){

	file, err := os.Open(ZabovHostsFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		d:= scanner.Text()
		DomainKill(d,ZabovHostsFile)
		incrementStats("Blacklist",1 )
		
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Println(err.Error())
	}



}






