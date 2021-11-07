package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	_ "embed"
)

//go:embed version
var version string

func main() {
	// read version at build
	fmt.Println(version)

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
		w.Write([]byte(getMetrics(config.ClientId, config.ClientSecret)))
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

func getToken(ClientId string, ClientSecret string) string {
	data := url.Values{
		"grant_type": {"password"},
		"username":   {ClientId},
		"password":   {ClientSecret},
	}

	client := &http.Client{}
	r, err := http.NewRequest("POST", "https://api-1.adax.no/client-api/auth/token", strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func getMetrics(ClientId string, ClientSecret string) string {
	// get token after JWT auth
	token := getToken(ClientId, ClientSecret)
	fmt.Println(token)

	var bearer = "Bearer " + token
	req, err := http.NewRequest("GET", "https://api-1.adax.no/client-api/rest/v1/control", nil)

	// add auth header
	req.Header.Add("Authorization", bearer)

	// perform request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
