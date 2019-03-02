// SEE LICENSE

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ghodss/yaml"
)

// IGDCfg is the configuration struct to hold 'golden config'
type IGDCfg struct {
	DNS    []string `yaml:"name_servers"`
	NTP    []string `yaml:"ntp_servers"`
	SYSLOG []string `yaml:"syslog_servers"`
}

func main() {
	// open the xlsx book of stuff
	xlsx, err := excelize.OpenFile("./goldenVars.xlsx")

	if err != nil {
		fmt.Println(err)
		return
	}

	// Now we're going to lazy load the IGDCfg struct
	cfg := &IGDCfg{}

	// Tracking variables for column indices
	DNSColumn := 0
	NTPColumn := 0
	SYSLOGColumn := 0

	for rowcount, row := range xlsx.GetRows("westEurope") {

		for colcount, colCell := range row {

			if rowcount == 0 {
				switch colCell {
				case "DNS":
					DNSColumn = colcount
				case "NTP":
					NTPColumn = colcount
				case "SYSLOG":
					SYSLOGColumn = colcount
				}
			}

			if rowcount > 0 {
				switch {
				case colcount == DNSColumn:
					if colCell != "" {
						cfg.DNS = append(cfg.DNS, colCell)
					}
				case colcount == NTPColumn:
					if colCell != "" {
						cfg.NTP = append(cfg.NTP, colCell)
					}
				case colcount == SYSLOGColumn:
					if colCell != "" {
						cfg.SYSLOG = append(cfg.SYSLOG, colCell)
					}
				}
			}
		}
	}

	// YAML Print to CLI Test
	cfgy, err := yaml.Marshal(&cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("---\n%s", string(cfgy))

	// Write file to group_vars
	f, err := os.OpenFile("./group_vars/westEurope/goldenvars.yaml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("---\n")); err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(cfgy)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
