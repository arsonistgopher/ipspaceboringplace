// Invoke with: xxxnyyynz file_output.yml
// xxx = custID
// yyy = orderID
// z = siteID

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ghodss/yaml"
)

// VARSDUMP main
type VARSDUMP struct {
	SiteID          string `json:"SiteID" yaml:"SiteID"`
	CUIDSID         string `json:"CUIDSID" yaml:"CUIDSID"`
	CUID            string `json:"CUID" yaml:"CUID"`
	SID             string `json:"SID" yaml:"SID"`
	SiteRD          string `json:"SiteRD" yaml:"SiteRD"`
	SiteRT          string `json:"SiteRT" yaml:"SiteRT"`
	LoopbackAddress string `json:"LoopbackAddress" yaml:"LoopbackAddress"`
	LoopbackUnit    string `json:"LoopbackUnit" yaml:"LoopbackUnit"`
	PEPort          string `json:"PEPort" yaml:"PEPort"`
	ASN             string `json:"ASN" yaml:"ASN"`
	PEAddress       string `json:"PEAddress" yaml:"PEAddress"`
	CPEAddress      string `json:"CPEAddress" yaml:"CPEAddress"`
}

func main() {

	cuidsid := os.Args[1]
	preURL := "http://localhost:1323/vars/%s"
	url := fmt.Sprintf(preURL, cuidsid)

	Client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, getErr := Client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	vars := VARSDUMP{}
	jsonErr := json.Unmarshal(body, &vars)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	YAMLCfg, err := yaml.Marshal(&vars)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(YAMLCfg))

	// Write files
	fileOutput := os.Args[2]

	// Write file to group_vars
	f1, err := os.OpenFile(fileOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f1.Write([]byte("---\n")); err != nil {
		log.Fatal(err)
	}
	if _, err := f1.Write([]byte(YAMLCfg)); err != nil {
		log.Fatal(err)
	}
	if err := f1.Close(); err != nil {
		log.Fatal(err)
	}
}
