// 1. Serve variables from DB
// 2. For POST accept customer:serviceid and return set of variables
// 3. Serve POST / GET / DELETE and return JSON

package main

import (
	"fmt"
	"net"
	"strings"

	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

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

// CFGVARS main
type CFGVARS struct {
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

//----------
// Handlers
//----------

func createVars(c echo.Context) error {

	VARS := &CFGVARS{}
	if err := c.Bind(VARS); err != nil {
		return err
	}

	cuidsidsiteSplit := strings.Split(VARS.CUIDSID, "n")
	VARS.CUID, VARS.SID, VARS.SiteID = strings.ToLower(cuidsidsiteSplit[0]), strings.ToLower(cuidsidsiteSplit[1]), strings.ToLower(cuidsidsiteSplit[2])

	// Get one of everything we need.
	var rts SiteRTRDs
	var loop LoopBacks
	var siteportsA SitePortsA
	var siteportsB SitePortsB
	var asn L3VPNASN
	var transit Transit
	// These are specific per site
	var loopunitA Sitelo0UnitsA
	var loopunitB Sitelo0UnitsB

	db, err := gorm.Open("sqlite3", "automationvars.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Check if it exists already in the DB by using check entry as a dirty check against LoopBacks
	var checks []LoopBacks
	if err := db.Where(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID}).Find(&checks).Error; err == nil {
		for _, check := range checks {
			if check.CUID == VARS.CUID && check.SID == VARS.SID {
				return c.NoContent(http.StatusUnprocessableEntity)
			}
		}

	}

	// Begin creation

	db.First(&rts)

	// Check site ID from request
	switch VARS.SiteID {
	case "a":
		if err := db.Where("Consumed = ?", 0).First(&loopunitA).Error; err == nil {
			db.Model(&loopunitA).Update(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID, "consumed": true})
			VARS.LoopbackUnit = loopunitA.Unit
		} else {
			fmt.Println("Loopback Unit exhaustion error")
			// We silence the error below here because it needs to go to a log ideally, but not litter use space
			deleteRowsByCUIDSID(VARS.CUID, VARS.SID)
			return c.NoContent(http.StatusUnprocessableEntity)
		}
	case "b":
		if err := db.Where("Consumed = ?", 0).First(&loopunitB).Error; err == nil {
			db.Model(&loopunitB).Update(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID, "consumed": true})
			VARS.LoopbackUnit = loopunitB.Unit
		} else {
			fmt.Println("Loopback Unit exhaustion error")
			// We silence the error below here because it needs to go to a log ideally, but not litter use space
			deleteRowsByCUIDSID(VARS.CUID, VARS.SID)
			return c.NoContent(http.StatusUnprocessableEntity)
		}
	}

	if err := db.Where("Consumed = ?", 0).First(&loop).Error; err == nil {
		db.Model(&loop).Update(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID, "consumed": true})
		VARS.LoopbackAddress = loop.IP.String()
	} else {
		fmt.Println("Loopback exhaustion error")
		// We silence the error below here because it needs to go to a log ideally, but not litter use space
		deleteRowsByCUIDSID(VARS.CUID, VARS.SID)
		return c.NoContent(http.StatusUnprocessableEntity)
	}

	// Check site ID from request
	switch VARS.SiteID {
	case "a":
		if err := db.Where("Consumed = ?", 0).First(&siteportsA).Error; err == nil {
			db.Model(&siteportsA).Update(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID, "consumed": true})
			VARS.PEPort = siteportsA.Port
		} else {
			fmt.Printf("Port exhaustion error: %s\n", err)
			// We silence the error below here because it needs to go to a log ideally, but not litter use space
			deleteRowsByCUIDSID(VARS.CUID, VARS.SID)
			return c.NoContent(http.StatusUnprocessableEntity)

		}
	case "b":
		if err := db.Where("Consumed = ?", 0).First(&siteportsB).Error; err == nil {
			db.Model(&siteportsB).Update(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID, "consumed": true})
			VARS.PEPort = siteportsB.Port
		} else {
			fmt.Printf("Port exhaustion error: %s\n", err)
			// We silence the error below here because it needs to go to a log ideally, but not litter use space
			deleteRowsByCUIDSID(VARS.CUID, VARS.SID)
			return c.NoContent(http.StatusUnprocessableEntity)

		}
	}

	if err := db.Where("Consumed = ?", 0).First(&asn).Error; err == nil {
		db.Model(&asn).Update(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID, "consumed": true})
		VARS.ASN = asn.ASN
	} else {
		fmt.Printf("ASN exhaustion error: %s\n", err)
		// We silence the error below here because it needs to go to a log ideally, but not litter use space
		deleteRowsByCUIDSID(VARS.CUID, VARS.SID)
		return c.NoContent(http.StatusUnprocessableEntity)
	}

	if err := db.Where("Consumed = ?", 0).First(&transit).Error; err == nil {
		db.Model(&transit).Update(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID, "consumed": true})
		VARS.PEAddress = transit.PE.String()
		VARS.CPEAddress = transit.CPE.String()
	} else {
		fmt.Printf("Transit Prefix exhaustion error: %s\n", err)
		// We silence the error below here because it needs to go to a log ideally, but not litter use space
		deleteRowsByCUIDSID(VARS.CUID, VARS.SID)
		return c.NoContent(http.StatusUnprocessableEntity)
	}

	VARS.SiteRT = rts.SiteRT
	VARS.SiteRD = rts.SiteRD

	fmt.Print(VARS)
	return c.JSON(http.StatusCreated, VARS)
}

// Get variables for cuid/sid
func getVars(c echo.Context) error {
	var loopbacks []LoopBacks
	var siteportsA []SitePortsA
	var siteportsB []SitePortsB
	var loopunitsA []Sitelo0UnitsA
	var loopunitsB []Sitelo0UnitsB
	var asns []L3VPNASN
	var transits []Transit
	var rts SiteRTRDs

	VARS := &CFGVARS{}

	cuidsidsite := c.Param("cuidsid")
	cuidsidsiteSplit := strings.Split(cuidsidsite, "n")
	VARS.CUIDSID = cuidsidsite
	VARS.CUID, VARS.SID, VARS.SiteID = strings.ToLower(cuidsidsiteSplit[0]), strings.ToLower(cuidsidsiteSplit[1]), strings.ToLower(cuidsidsiteSplit[2])

	// Get all other vars
	db, err := gorm.Open("sqlite3", "automationvars.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	if err := db.Where(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID}).Find(&loopbacks).Error; err == nil {
		for _, loopback := range loopbacks {
			VARS.LoopbackAddress = loopback.IP.String()
		}
	} else {
		return err
	}

	// Check site ID from request
	switch VARS.SiteID {
	case "a":

		if err := db.Where(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID}).Find(&loopunitsA).Error; err == nil {
			for _, loopaunit := range loopunitsA {
				VARS.LoopbackUnit = loopaunit.Unit
			}

		} else {
			return err
		}
	case "b":
		if err := db.Where(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID}).Find(&loopunitsB).Error; err == nil {
			for _, loopbunit := range loopunitsB {
				VARS.LoopbackUnit = loopbunit.Unit
			}

		} else {
			return err
		}
	}

	// Check site ID from request
	switch VARS.SiteID {
	case "a":
		if err := db.Where(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID}).Find(&siteportsA).Error; err == nil {
			for _, siteport := range siteportsA {
				VARS.PEPort = siteport.Port
			}

		} else {
			return err
		}
	case "b":
		if err := db.Where(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID}).Find(&siteportsB).Error; err == nil {
			for _, siteport := range siteportsB {
				VARS.PEPort = siteport.Port
			}

		} else {
			return err
		}
	}

	if err := db.Where(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID}).Find(&asns).Error; err == nil {
		for _, asn := range asns {
			VARS.ASN = asn.ASN
		}
	} else {
		return err
	}

	if err := db.Where(map[string]interface{}{"c_uid": VARS.CUID, "s_id": VARS.SID}).Find(&transits).Error; err == nil {
		for _, transit := range transits {
			prefixSize, _ := transit.Mask.Size()
			VARS.PEAddress = fmt.Sprintf("%s/%v", transit.PE.String(), prefixSize)
			VARS.CPEAddress = fmt.Sprintf("%s/%v", transit.CPE.String(), prefixSize)
		}
	} else {
		return err
	}

	db.First(&rts)

	VARS.SiteRT = rts.SiteRT
	VARS.SiteRD = rts.SiteRD

	return c.JSON(http.StatusOK, VARS)
}

// Deletion func from REST
func deleteVars(c echo.Context) error {
	cuidsid := c.Param("cuidsid")
	cuidsidSplit := strings.Split(cuidsid, "n")
	cuid, sid := cuidsidSplit[0], cuidsidSplit[1]

	err := deleteRowsByCUIDSID(cuid, sid)
	if err != nil {
		fmt.Print(err)
	}
	// Delete
	return c.NoContent(http.StatusNoContent)
}

// Function deletes stuff out of the database
func deleteRowsByCUIDSID(cuid string, sid string) error {
	var loopbacks []LoopBacks
	var siteportsA []SitePortsA
	var siteportsB []SitePortsB
	var loopunitsA []Sitelo0UnitsA
	var loopunitsB []Sitelo0UnitsB
	var asns []L3VPNASN
	var transits []Transit

	db, err := gorm.Open("sqlite3", "automationvars.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	if err := db.Where(map[string]interface{}{"c_uid": cuid, "s_id": sid}).Find(&loopunitsA).Error; err == nil {
		for _, loopunit := range loopunitsA {
			db.Model(&loopunit).Update(map[string]interface{}{"c_uid": "0", "s_id": "0", "consumed": false})
		}
	} else {
		return err
	}

	if err := db.Where(map[string]interface{}{"c_uid": cuid, "s_id": sid}).Find(&loopunitsB).Error; err == nil {
		for _, loopunit := range loopunitsB {
			db.Model(&loopunit).Update(map[string]interface{}{"c_uid": "0", "s_id": "0", "consumed": false})
		}
	} else {
		return err
	}

	if err := db.Where(map[string]interface{}{"c_uid": cuid, "s_id": sid}).Find(&loopbacks).Error; err == nil {
		for _, loopback := range loopbacks {
			db.Model(&loopback).Update(map[string]interface{}{"c_uid": "0", "s_id": "0", "consumed": false})
		}
	} else {
		return err
	}

	if err := db.Where(map[string]interface{}{"c_uid": cuid, "s_id": sid}).Find(&siteportsA).Error; err == nil {
		for _, siteport := range siteportsA {
			db.Model(&siteport).Update(map[string]interface{}{"c_uid": "0", "s_id": "0", "consumed": false})
		}
	} else {
		return err
	}

	if err := db.Where(map[string]interface{}{"c_uid": cuid, "s_id": sid}).Find(&siteportsB).Error; err == nil {
		for _, siteport := range siteportsB {
			db.Model(&siteport).Update(map[string]interface{}{"c_uid": "0", "s_id": "0", "consumed": false})
		}
	} else {
		return err
	}

	if err := db.Where(map[string]interface{}{"c_uid": cuid, "s_id": sid}).Find(&asns).Error; err == nil {
		for _, asn := range asns {
			db.Model(&asn).Update(map[string]interface{}{"c_uid": "0", "s_id": "0", "consumed": false})
		}
	} else {
		return err
	}

	if err := db.Where(map[string]interface{}{"c_uid": cuid, "s_id": sid}).Find(&transits).Error; err == nil {
		for _, transit := range transits {
			db.Model(&transit).Update(map[string]interface{}{"c_uid": "0", "s_id": "0", "consumed": false})
		}
	} else {
		return err
	}

	return nil
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/vars", createVars)
	e.GET("/vars/:cuidsid", getVars)
	e.DELETE("/vars/:cuidsid", deleteVars)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
