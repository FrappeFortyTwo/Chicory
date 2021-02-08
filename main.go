package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
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

	// regex to fetch video urls
	//re := regexp.MustCompile(`{"itag"+.+\d+.+}`)
	re := regexp.MustCompile(`"adaptiveFormats"+.+\]\},"playerAds"`)

	// get video source & meta in json format
	vs := strings.Replace(string(re.FindAll(body, -1)[0]), "\"adaptiveFormats\":[", "", 1)
	vs = strings.Replace(vs, "]},\"playerAds\"", "", 1)

	println(vs)

	// dump contents into json file
	// err = ioutil.WriteFile("temp.json", vs[0], 0777)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
