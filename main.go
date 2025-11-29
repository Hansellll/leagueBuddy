package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	//User input
	in := bufio.NewReader(os.Stdin)

	var word string
	fmt.Print("Please enter your username: \n")
	if _, err := fmt.Fscan(in, &word); err != nil {
		fmt.Println("scan error:", err)
		return
	}
	fmt.Println("You typed: ", word)

	if len(os.Args) < 2 {
		log.Fatal("usage: go run main.go <full-url>")
	}
	url := os.Args[1]
	apiKey := os.Getenv("RIOT_API_KEY")
	if apiKey == "" {
		log.Fatal("RIOT_API_KEY env var is not set")
	}
	if err := api(url, apiKey); err != nil {
		log.Fatal(err)
	}
}

func api(url, apiKey string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Riot-Token", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.Status)
	fmt.Println(string(b))
	return nil
}
