package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0"
	BaseEndpoint = "https://www.instagram.com/api/v1/users/web_profile_info/?username=%s"
	IGAppID   = "936619743392459"
)

// TODO "I want it profile.ID not profile.Data.User.ID"
type APIResponse struct {
	Data struct {
		User struct {
			ID             string `json:"id"`
			FullName       string `json:"full_name"`
			Username       string `json:"username"`
			Biography      string `json:"biography"`
			BioLink        string `json:"external_url"`
			IsPrivate      bool   `json:"is_private"`
			IsVerified     bool   `json:"is_verified"`
			FollowersCount int    `json:"edge_followed_by.count"`
			FollowingCount int    `json:"edge_follow.count"`
		} `json:"user"`
	} `json:"data"`
}

func LookupProfile(username string) (APIResponse, error) {
	url := fmt.Sprintf(BaseEndpoint, username)

	//fmt.Println(url)

	// TODO create separated fetch func
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return APIResponse{}, errors.New("error creating new request")
	}

	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("X-IG-App-ID", IGAppID)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return APIResponse{}, errors.New("error sending request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return APIResponse{}, errors.New("user not found")
		}
		return APIResponse{}, fmt.Errorf("error %s", http.StatusText(response.StatusCode))
	}

	// TODO save the original response (not truncated) to a file `output/{username}.json`
	var apiResponse APIResponse
	if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
		return APIResponse{}, fmt.Errorf("error parsing json response - %v", err)
	}

	return apiResponse, nil
}

func main() {

    // TODO make dynamic username using prompt/argument(?)
	username := "zuck"

	lookup, err := LookupProfile(username)
	if err != nil {
		fmt.Printf("\033[31m%s\033[0m\n", err)
		return
	}

	// TODO prettify output
	profile := lookup.Data.User
	fmt.Printf("Profile: %+v\n", profile)
}
