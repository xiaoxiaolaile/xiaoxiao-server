package main

import (
	"fmt"
	"log"
)
import "github.com/boltdb/bolt"

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("sillyGirl.cache", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//db.View(func(tx *bolt.Tx) error {
	//	// Assume bucket exists and has keys
	//	b := tx.Bucket([]byte("plugins"))
	//	b.ForEach(func(k, v []byte) error {
	//		fmt.Printf("key=%s, value=%s\n", k, v)
	//		fmt.Println("=================================")
	//		return nil
	//	})
	//	return nil
	//})

	db.View(func(tx *bolt.Tx) error {

		tx.ForEach(func(name []byte, b *bolt.Bucket) error {

			fmt.Printf("name=%s, \n", name)

			return nil
		})

		return nil
	})

	//db.View(func(tx *bolt.Tx) error {
	//	b := tx.Bucket([]byte("plugins"))
	//	v := b.Get([]byte("ccc4bbc0-3187-11ed-b60c-aaaa00117a5c"))
	//	fmt.Printf("The answer is: %s\n", v)
	//	return nil
	//})
}
