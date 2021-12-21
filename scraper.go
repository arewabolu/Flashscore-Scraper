package scraper

import (
	"io"
	"log"
	"net/http"
	"os"
)

func pingMan(weblink, meth string) {
	client := http.Client{}
	request, err := http.NewRequest(meth, weblink, nil)
	if err != nil {
		log.Println(err)
	}
	request.ContentLength
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	if response.Header.Get("Content-Type") == "text/html" {
		file, err := os.Create("euroleague.html")
		if err != nil {
			log.Fatalln(err)
		}
	}
	defer file.Close()

	// Copy response body to file
	newFile, err := io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	meth := os.Args[1]
	var weblink string

	switch meth {
	case "GET":
		http.ReadResponse()
	}
}
