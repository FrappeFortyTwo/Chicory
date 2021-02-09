package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Formats struct which contains array of formats
type Formats struct {
	Formats []Format `json:"formats"`
}

// Format struct which contains Format meta data
type Format struct {
	Itag             int        `json:"itag"`
	MimeType         string     `json:"mimeType"`
	Bitrate          int32      `json:"bitrate"`
	Width            int        `json:"width"`
	Height           int        `json:"height"`
	InitRange        InitRange  `json:"initRange"`
	IndexRange       IndexRange `json:"indexRange"`
	LastModified     string     `json:"lastModified"`
	ContentLength    string     `json:"contentLength"`
	Quality          string     `json:"quality"`
	Fps              int        `json:"fps"`
	QualityLabel     string     `json:"qualityLabel"`
	ProjectionType   string     `json:"projectionType"`
	AverageBitrate   int32      `json:"averageBitrate"`
	ApproxDurationMs string     `json:"approxDurationMS"`
	SignatureCipher  string     `json:"signatureCipher"`
}

// InitRange which contains it's Start and End
type InitRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// IndexRange contains it's Start and End
type IndexRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func main() {

	// define pointers to command-line arguments
	url := flag.String("url", "", "web address : (something like this) https://www.youtube.com/watch?v=RR8dqCCZ_IY")
	bulk := flag.Bool("bulk", false, "false  : single video\ntrue   : multiple videos i.e playlist")

	// parse command-line arguments
	flag.Parse()

	println("\n // ---------- Chicory Youtube Video Downloader ---------- // \n")
	println("* Fetching video from  : ", *url)
	println("* Bulk download option : ", *bulk)
	println()

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

	// regex to fetch video urls & meta data
	re := regexp.MustCompile(`"adaptiveFormats"+.+\]\},"playerAds"`)

	// process data into json format
	vs := strings.Replace(string(re.FindAll(body, -1)[0]), "\"adaptiveFormats\":", "{ \"formats\":", 1)
	vs = strings.Replace(vs, "},\"playerAds\"", "}", 1)

	// dump contents into file
	err = ioutil.WriteFile("temp.json", []byte(vs), 0777)
	if err != nil {
		log.Fatal(err)
	}

	// read file as json
	jsonFile, err := os.Open("temp.json")
	if err != nil {
		log.Fatalln(err)
	}

	defer jsonFile.Close()

	byteVal, _ := ioutil.ReadAll(jsonFile)

	// initialise formats ~ various formats the video at provided url is available in
	var formats Formats

	// unmarshal contents
	json.Unmarshal(byteVal, &formats)

	// iterate through every format and print respective meta data
	println("Option\t|\tItag\t|\tType\t\t|\tQuality\n")
	for i := 0; i < len(formats.Formats); i++ {

		tmpA := strings.Split(formats.Formats[i].MimeType, "; ")
		println(i, "\t|\t", formats.Formats[i].Itag, "\t|\t", tmpA[0], "\t|\t", formats.Formats[i].QualityLabel, "\n")
	}

}
