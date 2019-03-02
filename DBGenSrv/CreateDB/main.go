package main

import (
	"fmt"
	"net"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Variable pool we want served!
var siteRT = "target:64501:1"
var siteRD = "64501:1"
var sitelo0unitsA = []string{"2", "3", "4", "5"}
var sitelo0unitsB = []string{"2", "3", "4", "5"}
var loopbacks = []string{"192.0.192.1", "192.0.192.2", "192.0.192.3", "192.0.192.4", "192.0.192.5", "192.0.192.6", "192.0.192.7", "192.0.192.8"}
var transitprefixes = []string{"172.16.1.0/30", "172.16.1.4/30", "172.16.1.8/30", "172.16.1.12/30"}
var asns = []string{"64502", "64503", "64504", "64505", "64506"}
var siteportsA = []string{"ge-0/0/1", "ge-0/0/2", "ge-0/0/3", "ge-0/0/4", "ge-0/0/5"}
var siteportsB = []string{"ge-0/0/1", "ge-0/0/2", "ge-0/0/3", "ge-0/0/4", "ge-0/0/5"}

// Sitelo0UnitsA consumption
type Sitelo0UnitsA struct {
	gorm.Model
	CUID     string
	SID      string
	Consumed bool
	Unit     string
}

// Sitelo0UnitsB consumption
type Sitelo0UnitsB struct {
	gorm.Model
	CUID     string
	SID      string
	Consumed bool
	Unit     string
}

// SiteRTRDs consumption
type SiteRTRDs struct {
	gorm.Model
	SiteRT string
	SiteRD string
}

// LoopBacks consumption
type LoopBacks struct {
	gorm.Model
	CUID     string
	SID      string
	Consumed bool
	IP       net.IP
}

// SitePortsA consumption
type SitePortsA struct {
	gorm.Model
	CUID     string
	SID      string
	Consumed bool
	Port     string
}

// SitePortsB consumption
type SitePortsB struct {
	gorm.Model
	CUID     string
	SID      string
	Consumed bool
	Port     string
}

// L3VPNASN is used for tracking ASNs
type L3VPNASN struct {
	gorm.Model
	CUID     string
	SID      string
	Consumed bool
	ASN      string
}

// Transit is used to store /30 prefixes
type Transit struct {
	gorm.Model
	CUID     string
	SID      string
	Consumed bool
	Net      net.IP
	Mask     net.IPMask
	PE       net.IP
	CPE      net.IP
}

func main() {
	db, err := gorm.Open("sqlite3", "automationvars.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Transit{})
	db.AutoMigrate(&L3VPNASN{})
	db.AutoMigrate(&SitePortsA{})
	db.AutoMigrate(&SitePortsB{})
	db.AutoMigrate(&LoopBacks{})
	db.AutoMigrate(&Sitelo0UnitsA{})
	db.AutoMigrate(&Sitelo0UnitsB{})
	db.AutoMigrate(&SiteRTRDs{})

	// Create all DB items
	// Start with SiteRTs
	db.Create(&SiteRTRDs{SiteRT: siteRT, SiteRD: siteRD})

	// Next, create loopbacks
	for _, lo := range loopbacks {

		IP := net.ParseIP(lo)

		if err != nil {
			fmt.Print(err)
			panic(err)
		}
		db.Create(&LoopBacks{IP: IP})
	}

	// Next, create SiteA lo0 units
	for _, unit := range sitelo0unitsA {

		if err != nil {
			fmt.Print(err)
			panic(err)
		}
		db.Create(&Sitelo0UnitsA{Unit: unit})
	}

	// Next, create SiteB lo0 units
	for _, unit := range sitelo0unitsB {

		if err != nil {
			fmt.Print(err)
			panic(err)
		}
		db.Create(&Sitelo0UnitsB{Unit: unit})
	}

	// Next, create ASNs
	for _, asn := range asns {

		if err != nil {
			fmt.Print(err)
			panic(err)
		}
		db.Create(&L3VPNASN{ASN: asn})
	}

	// Next, create SiteAPorts
	for _, port := range siteportsA {

		if err != nil {
			fmt.Print(err)
			panic(err)
		}
		db.Create(&SitePortsA{Port: port})
	}

	// Next, create SiteBPorts
	for _, port := range siteportsB {

		if err != nil {
			fmt.Print(err)
			panic(err)
		}
		db.Create(&SitePortsB{Port: port})
	}

	// Next, create TransitPrefixes
	for _, prefix := range transitprefixes {

		if err != nil {
			fmt.Print(err)
			panic(err)
		}

		_, IPNet, err := net.ParseCIDR(prefix)
		if err != nil {
			fmt.Print(err)
			panic(err)
		}

		// Apply mask of prefix and get network address
		BoundaryCheckMask := IPNet.IP.Mask(IPNet.Mask)

		// Increment last octet to get end-point address
		netStr := IPNet.IP.String()
		IPNet.IP[3]++
		ep1Str := IPNet.IP.String()
		IPNet.IP[3]++
		ep2Str := IPNet.IP.String()

		netAddr := net.ParseIP(netStr)
		ep1 := net.ParseIP(ep1Str)
		ep2 := net.ParseIP(ep2Str)

		if BoundaryCheckMask.String() == ep1.Mask(IPNet.Mask).String() && BoundaryCheckMask.String() == ep2.Mask(IPNet.Mask).String() {
			db.Create(&Transit{Net: netAddr, Mask: IPNet.Mask, PE: ep1, CPE: ep2})
		} else {
			panic("Error with IPv4 incrementing")
		}
	}
}
