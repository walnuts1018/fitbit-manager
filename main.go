package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	urlPrefix := "https://api.fitbit.com/1/user/" + config.UserID
	//url := urlPrefix "/sleep/date/2023-09-15.json"
	url := urlPrefix + "/activities/heart/date/today/today/1sec/time/15:00/23:00.json"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	byteArray, _ := io.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
	fmt.Printf("%#v\n", resp.Header)
}

func refresh(config *config) {
	req, _ := http.NewRequest("GET", "https://api.fitbit.com/oauth2/token", nil)
	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

}
