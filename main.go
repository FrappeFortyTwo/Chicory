package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	// define pointers to command-line arguments
	url := flag.String("url", "", "web address : (something like this) https://www.youtube.com/watch?v=RR8dqCCZ_IY")
	bulk := flag.Bool("bulk", false, "false  : single video\ntrue   : multiple videos i.e playlist")

	// parse command-line arguments
	flag.Parse()

	println(*url)
	println(*bulk)

	// fetch source for url
	resp, err := http.Get(*url)
	if err != nil {
		log.Fatalln(err)
	}
	// close after usage
	defer resp.Body.Close()

	// read contents from url response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	}

	// dump contents into a file
	// err = ioutil.WriteFile("temp.txt", body, 0777)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	println(string(body))
}
