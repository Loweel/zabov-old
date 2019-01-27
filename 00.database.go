package main

import (
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

//MyZabovDB is the storage where we'll put domains to block
var MyZabovDB *bolt.DB

func init() {

	var err error

	os.RemoveAll("./db")

	os.MkdirAll("./db", 0755)

	MyZabovDB, err = bolt.Open("./db/zabov.db", 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("Problem Creating Zabov DB: ", err.Error())
	} else {
		fmt.Println("Zabov DB Created")

	}

	err = MyZabovDB.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists(zabovKbucket)
		if err != nil {
			fmt.Printf("could not create %s bucket: %v\n", zabovKbucket, err)
			return err
		}
		fmt.Println("Created bucket:", string(zabovKbucket))

		_, err = root.CreateBucketIfNotExists(zabovCbucket)
		if err != nil {
			fmt.Printf("could not create %s bucket: %v\n", zabovCbucket, err)
			return err
		}
		fmt.Println("Created bucket:", string(zabovCbucket))

		return nil
	})

}
