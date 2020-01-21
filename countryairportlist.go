package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// CountryDetail is a type which shows details of countries
type CountryDetail struct {
	ID            int    `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Continent     string `json:"continent"`
	WikipediaLink string `json:"wikipedia_link"`
	Keywords      string `json:"keywords"`
}

// AirportDetail is a type which shows details of the airports...
type AirportDetail struct {
	ID               int     `json:"id"`
	Ident            string  `json:"ident"`
	AType            string  `json:"aType"`
	Name             string  `json:"name"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Elevation        int     `json:"elevation"`
	Continent        string  `json:"continent"`
	IsoCountry       string  `json:"iso_country"`
	IsoRegion        string  `json:"iso_region"`
	Municipality     string  `json:"municipality"`
	ScheduledService string  `json:"scheduled_service"`
	GpsCode          string  `json:"gps_code"`
	LocalCode        string  `json:"local_code"`
	Runways          []struct {
		ID           int    `json:"id"`
		AirportRef   int    `json:"airport_ref"`
		AirportIdent string `json:"airport_ident"`
		LengthFt     int    `json:"length_ft"`
		WidthFt      int    `json:"width_ft"`
		Surface      string `json:"surface"`
		Lighted      int    `json:"lighted"`
		Closed       int    `json:"closed"`
		LeIdent      string `json:"le_ident"`
		HeIdent      string `json:"he_ident"`
	} `json:"runways"`
}

// CountryAirportSummary is a type which shows details of the airport which the country has...
type CountryAirportSummary struct {
	Code        string `json:"code"`
	Countryname string `json:"countryname"`
	AirportSum  []struct {
		Airportid        string `json:"airportid"`
		Airportname      string `json:"airportname"`
		Numberofrunaways int    `json:"numberofrunaways"`
	} `json:"airportsum"`
}

// AirportSum is a type which shows summary of airport details
type AirportSum struct {
	Airportid        string `json:"airportid"`
	Airportname      string `json:"airportname"`
	Numberofrunaways int    `json:"numberofrunaways"`
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func fetchCountryDetails() []CountryDetail {
	var countryDetailMessage []CountryDetail
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "http://country-api.info/countries", nil)
	if err != nil {
		log.Print(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
	} else {
		jsonResponse, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
		}
		defer resp.Body.Close()
		json.Unmarshal([]byte(jsonResponse), &countryDetailMessage)
	}
	return countryDetailMessage
}

// func fetchAirportDetails() []AirportDetail {
// 	var airportDetailMessage []AirportDetail
// 	client := &http.Client{Timeout: 300 * time.Second}
// 	req, err := http.NewRequest("GET", "http://airport-api.info/airports", nil)
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Close = true
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		jsonResponse, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			//log.Fatal(err)
// 			log.Print(err)
// 		}
// 		defer resp.Body.Close()
// 		json.Unmarshal([]byte(jsonResponse), &airportDetailMessage)

// 	}
// 	return airportDetailMessage
// }

func fetchAirportDetails() []AirportDetail {
	var airportDetailMessage []AirportDetail
	jsonResponse, err := ioutil.ReadFile("airportresponse.json")
	if err != nil {
		log.Print(err)
	}
	json.Unmarshal([]byte(jsonResponse), &airportDetailMessage)
	return airportDetailMessage
}

func getAllData(w http.ResponseWriter, r *http.Request) {
	countryDetailMessage := fetchCountryDetails()
	airportDetailMessage := fetchAirportDetails()

	//merge all data

	var summaryMessage [247]CountryAirportSummary

	for i := 0; i < len(countryDetailMessage); i++ {
		summaryMessage[i].Code = countryDetailMessage[i].Code
		summaryMessage[i].Countryname = countryDetailMessage[i].Name
	}

	for i := 0; i < len(summaryMessage); i++ {
		for k := 0; k < len(airportDetailMessage); k++ {
			if summaryMessage[i].Code == airportDetailMessage[k].IsoCountry {
				var temp AirportSum
				temp.Airportid = airportDetailMessage[k].Ident
				temp.Airportname = airportDetailMessage[k].Name
				temp.Numberofrunaways = len(airportDetailMessage[k].Runways) + 1
				summaryMessage[i].AirportSum = append(summaryMessage[i].AirportSum, temp)
			}
		}
	}

	json.NewEncoder(w).Encode(summaryMessage)
}

// example: call /countryairportsummary/2 or /countryairportsummary/runwayminimum=2
//the function provides the required query parameter for runwayminimum
func returnRequiredQueryFunc(w http.ResponseWriter, r *http.Request) {
	countryDetailMessage := fetchCountryDetails()
	airportDetailMessage := fetchAirportDetails()

	//merge all data

	var summaryMessage [247]CountryAirportSummary

	for i := 0; i < len(countryDetailMessage); i++ {
		summaryMessage[i].Code = countryDetailMessage[i].Code
		summaryMessage[i].Countryname = countryDetailMessage[i].Name
	}

	for i := 0; i < len(summaryMessage); i++ {
		for k := 0; k < len(airportDetailMessage); k++ {
			if summaryMessage[i].Code == airportDetailMessage[k].IsoCountry {
				var temp AirportSum
				temp.Airportid = airportDetailMessage[k].Ident
				temp.Airportname = airportDetailMessage[k].Name
				temp.Numberofrunaways = len(airportDetailMessage[k].Runways) + 1
				summaryMessage[i].AirportSum = append(summaryMessage[i].AirportSum, temp)
			}
		}
	}

	// filter minimum runway of airports
	var tempAirportMinimumNumberOfRunAway []CountryAirportSummary

	minimumNumberOfRunAways := mux.Vars(r)["runwayminimum"]
	if minimumNumberOfRunAways == "" {
		minimumNumberOfRunAways = "0"
	}

	tempMinimumNumberOfRunAways, err := strconv.Atoi(minimumNumberOfRunAways)
	if err != nil {
		log.Print(err)
	}
	for i := 0; i < len(summaryMessage); i++ {
		var tempFilterCountryAirportSum CountryAirportSummary

		for k := 0; k < len(summaryMessage[i].AirportSum); k++ {
			if summaryMessage[i].AirportSum[k].Numberofrunaways >= tempMinimumNumberOfRunAways {
				tempFilterCountryAirportSum.Code = summaryMessage[i].Code
				tempFilterCountryAirportSum.Countryname = summaryMessage[i].Countryname
				tempAirportInfo := summaryMessage[i].AirportSum[k]
				tempFilterCountryAirportSum.AirportSum = append(tempFilterCountryAirportSum.AirportSum, tempAirportInfo)
			}
		}
		if tempFilterCountryAirportSum.Code != "" {
			tempAirportMinimumNumberOfRunAway = append(tempAirportMinimumNumberOfRunAway, tempFilterCountryAirportSum)
		}

	}
	json.NewEncoder(w).Encode(tempAirportMinimumNumberOfRunAway)

}

// HealthCheckHandler is a function that shows status of the service
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	log.Print("the service is working...")

	router.HandleFunc("/", homeLink)
	//router.HandleFunc("/countryairportsummary", getAllData).Methods("GET")
	router.HandleFunc("/countryairportsummary/{runwayminimum}", returnRequiredQueryFunc).Methods("GET")

	router.Path("/countryairportsummary").Queries("runwayminimum", "{runwayminimum}").HandlerFunc(returnRequiredQueryFunc).Name("returnRequiredQueryFunc").Methods("GET")
	router.Path("/countryairportsummary").HandlerFunc(getAllData).Methods("GET")
	router.HandleFunc("/health", HealthCheckHandler)

	log.Fatal(http.ListenAndServe(":8000", router))

}
