package main

/*
//https://github.com/Hansellll/leagueBuddy
@date: 11-29-25
@author: Hanselll
*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// JSON Struct for handling ACCOUNT-V1 repsonse
type Account struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

func main() {

	// Create bufio reader taking Stdin as input
	userin := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: \n")
	username, err := userin.ReadString('\n')
	if err != nil {
		log.Fatal("error reading username: ", err)
	}

	// Format for url
	username = strings.TrimSpace(username)
	// Get tagline
	fmt.Print("Enter tagline: \n")
	tagline, err := userin.ReadString('\n')
	if err != nil {
		log.Fatal("error reading your tagline: ", err)
	}
	// Format for url
	tagline = strings.TrimSpace(tagline)

	// Retrieve and verify api key, must be stored in env. variable RIOT_API_KEY
	apiKey := os.Getenv("RIOT_API_KEY")
	if apiKey == "" {
		log.Fatal("RIOT_API_KEY env var is not set")
	}
	// Clean up whitespace/quotes if any
	apiKey = strings.Trim(apiKey, " \t\r\n\"")

	fmt.Printf(
		"Your username is %q and your tagline is %q \n",
		username,
		tagline,
	)

	// Get account puuid via ACCOUNT-V1 endpoint
	account, err := getPuuid(username, tagline, apiKey)
	if err != nil {
		log.Fatal("getPuuid error: ", err)
	}

	fmt.Printf("Found account: %s#%s\n", account.GameName, account.TagLine)
	fmt.Println("PUUID: ", account.PUUID)

	// Get 5 last match IDs
	matches, err := getRecentMatches(account.PUUID, apiKey)
	if err != nil {
		log.Fatal("getRecentMatches error: ", err)
	}

	// getRecentMatches test method
	fmt.Println("Last 5 match IDs:")
	for _, m := range matches {
		fmt.Println(m)
	}

	//Use match-v5 api to extract data on previous 5 matches
	fmt.Println("Lets print some info on those match's! ")
	for _, m := range matches {
		fullurl := fmt.Sprintf(
			"https://americas.api.riotgames.com/lol/match/v5/matches/%s",
			url.PathEscape(m))
		fmt.Println("Match: ", m)
		if err := api(fullurl, apiKey); err != nil {
			log.Println("error fetching match", m, ":", err)
		}
		fmt.Println()
	}
}
func getPuuid(username, tagline, apiKey string) (*Account, error) {
	baseurl := "https://americas.api.riotgames.com"
	endpoint := fmt.Sprintf(
		"/riot/account/v1/accounts/by-riot-id/%s/%s",
		url.PathEscape(username),
		url.PathEscape(tagline),
	)
	fullurl := baseurl + endpoint

	req, err := http.NewRequest("GET", fullurl, nil)
	if err != nil {
		return nil, err
	}

	// API requires apikey be stored in header for ACCOUNT-V1
	req.Header.Set("X-Riot-Token", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("account API returned %s", resp.Status)
	}

	var acc Account
	if err := json.NewDecoder(resp.Body).Decode(&acc); err != nil {
		return nil, err
	}

	return &acc, nil
}

func getRecentMatches(
	puuid,
	apiKey string,
) ([]string, error) {
	baseurl := "https://americas.api.riotgames.com"

	// Adjust start/count as you like
	endpoint := fmt.Sprintf(
		"/lol/match/v5/matches/by-puuid/%s/ids?start=0&count=5",
		url.QueryEscape(puuid),
	)
	fullurl := baseurl + endpoint

	req, err := http.NewRequest("GET", fullurl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Riot-Token", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("match API returned %s", resp.Status)
	}

	// match-v5 /ids returns a JSON array of strings
	var matches []string
	if err := json.NewDecoder(resp.Body).Decode(&matches); err != nil {
		return nil, err
	}

	return matches, nil
}

func api(url, apiKey string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	//API requires apikey be stored in header for ACCOUNT-V1
	req.Header.Set("X-Riot-Token", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//Print response body
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.Status)
	fmt.Println(string(b))
	return nil
}
