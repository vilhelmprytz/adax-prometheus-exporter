package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"

	_ "embed"
)

//go:embed version
var version string

func main() {
	log.Println("Running adax-prometheus-exporter v" + version)

	// get path to conifg file if path is specified
	var configPath string
	for i, s := range os.Args {
		if s == "-c" || s == "--config" {
			if i+1 < len(os.Args) {
				configPath = os.Args[i+1]
			}
		}
	}

	if configPath == "" {
		log.Fatal("Missing config parameter, -c or --config")
	}

	// read config file
	config := readConfig(configPath)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		resp, err := getMetrics(config.ClientId, config.ClientSecret)

		if err != nil {
			log.Println("ERROR", err)
			w.WriteHeader(http.StatusInternalServerError)
			resp = "Internal server error"
		}

		w.Write([]byte(resp))
	})
	log.Println("Listening on port", config.Port)
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(config.Port), nil))
}

type Config struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Port         int    `yaml:"port"`
}

func readConfig(path string) Config {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func getToken(ClientId string, ClientSecret string) (string, error) {
	response, err := http.PostForm("https://api-1.adax.no/client-api/auth/token", url.Values{
		"grant_type": {"password"},
		"username":   {ClientId},
		"password":   {ClientSecret}})

	//okay, moving on...
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("error getting token: %s", response.Status)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	type Token struct {
		AccessToken string `json:"access_token"`
	}

	var token Token
	json.Unmarshal([]byte(string(body)), &token)

	return token.AccessToken, nil
}

type Home struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Room struct {
	Id                int     `json:"id"`
	HomeId            int     `json:"homeId"`
	Name              string  `json:"name"`
	HeatingEnabled    bool    `json:"heatingEnabled"`
	TargetTemperature float64 `json:"targetTemperature"`
	Temperature       float64 `json:"temperature"`
}

type Data struct {
	Homes []Home `json:"homes"`
	Rooms []Room `json:"rooms"`
}

func getData(Token string) (Data, error) {
	req, err := http.NewRequest("GET", "https://api-1.adax.no/client-api/rest/v1/content/", nil)

	if err != nil {
		return Data{}, err
	}

	// add auth header
	req.Header.Add("Authorization", "Bearer "+Token)

	// perform request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return Data{}, err
	}

	if resp.StatusCode != 200 {
		return Data{}, fmt.Errorf("error getting homes: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Data{}, err
	}

	var data Data
	json.Unmarshal([]byte(string(body)), &data)

	return data, nil
}

func getMetrics(ClientId string, ClientSecret string) (string, error) {
	// get token after JWT auth
	token, err := getToken(ClientId, ClientSecret)

	if err != nil {
		return "", err
	}

	// get data
	data, err := getData(token)

	if err != nil {
		return "", err
	}

	var metrics string

	for _, home := range data.Homes {
		for _, room := range data.Rooms {
			if room.HomeId == home.Id {
				metrics += fmt.Sprintf("room_temperature{home=\"%s\",room=\"%s\"} %f\n", home.Name, room.Name, room.Temperature/100)
				metrics += fmt.Sprintf("room_target_temperature{home=\"%s\",room=\"%s\"} %f\n", home.Name, room.Name, room.TargetTemperature/100)
			}
		}
	}

	return metrics, nil
}
