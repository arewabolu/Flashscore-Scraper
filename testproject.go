package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	client := &http.Client{}
	// Make HTTP GET request
	request, err := http.NewRequest("GET", "https://www.euroleague.net/main/standings", nil)
	//www.eurocupbasketball.com/eurocup/games/standings

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	file, err := os.CreateTemp("", "league.html")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// Copy data from the response to standard output
	/*, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}*/

	n, err := io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	//var reader = bufio.NewReader(os.Stdout)
	//message, _ := reader.ReadString('\n')

	//log.Println("Number of bytes copied to STDOUT:", n)

	//fmt.Println(message)

}
