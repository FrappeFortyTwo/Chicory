package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
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
	URL              string     `json:"url"`
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
	checkErr(err, "Unable to make http request")
	defer resp.Body.Close()

	// read contents from url response
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err, "Unable to read response body")

	// regex to fetch video urls & meta data
	re := regexp.MustCompile(`"adaptiveFormats"+.+\]\},"playerAds"`)

	// process data into json format
	vs := strings.Replace(string(re.FindAll(body, -1)[0]), "\"adaptiveFormats\":", "{ \"formats\":", 1)
	vs = strings.Replace(vs, "},\"playerAds\"", "}", 1)

	// dump contents into file
	err = ioutil.WriteFile("temp.json", []byte(vs), 0777)
	checkErr(err, "Unable to write to temp.json")

	// read file as json
	jsonFile, err := os.Open("temp.json")
	checkErr(err, "Unable to read temp.json")
	defer jsonFile.Close()

	byteVal, _ := ioutil.ReadAll(jsonFile)

	// initialise formats ~ various formats the video at provided url is available in
	var formats Formats

	// unmarshal contents
	json.Unmarshal(byteVal, &formats)

	// iterate through every format and print respective meta data
	println("Option\t|\tItag\t|\tType\t\t|\tQuality\n")
	for i := 0; i < len(formats.Formats); i++ {

		// if meme type contains audio ~ break
		if strings.Contains(formats.Formats[i].MimeType, "audio") {
			break
		}

		// display options to download video
		tmpA := strings.Split(formats.Formats[i].MimeType, "; ")
		println(i, "\t|\t", formats.Formats[i].Itag, "\t|\t", tmpA[0], "\t|\t", formats.Formats[i].QualityLabel, "\n")
	}

	// input option to download video
	print("Enter Option: ")
	var input int
	fmt.Scanln(&input)

	// making channel of type string
	c := make(chan bool)

	// fetch video from url
	println("\nDownloading Video ...")
	go fetchFile(c, formats.Formats[input].URL, "temp-video")

	// fetch audio from url
	println("Downloading Audio ...")
	go fetchFile(c, formats.Formats[input+1].URL, "temp-audio")

	if <-c && <-c {
		println("\nMerging Files ...")
	}

	println("\nDownload Complete !")

}

func fetchFile(c chan bool, url string, title string) {

	// process url ~ replace u2600 to &
	url = strings.Replace(url, "u0026", "&", -1)

	// make http request &
	resp, err := http.Get(url)
	checkErr(err, "Unable to make http request")
	defer resp.Body.Close()

	// save response to file
	out, err := os.Create(title)
	checkErr(err, "Unable to create asset")
	defer out.Close()
	io.Copy(out, resp.Body)

	// return data to indicate task completion
	c <- true
}

func mergeFiles() {

}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(err, msg)
	}
}
