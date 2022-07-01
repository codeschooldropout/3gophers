package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/codeschooldropout/3gophers/logit"
)

const (
	VERSION = "0.0.1"
)

// TODO Redo struct, break into smaller structs

type Alert struct {
	Call            string  `json:"call"`      // entry, exit, error
	Type            string  `json:"type"`      // long, short, take
	Exchange        string  `json:"exchange"`  // 1. kucoin 2. kraken 3. coinbase - these are the ones i want to start with if possible.
	Base            string  `json:"base"`      // BTC, ETH, etc
	Quote           string  `json:"quote"`     // USD, EUR, etc
	Price           float64 `json:"price"`     // buy high, sell low right?
	Timeframe       string  `json:"timeframe"` // 1m 5m 15m 1h 4h 8h 1d etc
	PNL             float64 `json:"pnl"`       // profit/loss from exit signals
	Bars            int     `json:"bars"`      // Number of bars from entry signal to exit signal (bars x timeframe to identify corresponding signals)
	StopLoss        float64 `json:"sl"`        // Stop loss price to set on entry signal
	StopLossPercent float64 `json:"slp"`       // Stop loss percent to set on entry signal

}

func handleRP(w http.ResponseWriter, r *http.Request) {
	// make regex to match numbers and certain characters
	KeepNumbersReg, err := regexp.Compile("[^0-9-.]")
	if err != nil {
		logit.Log.Fatal(err)
	}

	// Read the body of the request
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logit.Log.Fatal(err)
	}

	// Log Initial Request Data both to raw file unchanged, to log file with new lines removed, and console
	logit.Raw.Printf("%s\n%s", r.URL.Path, reqBody)
	// logit.Log.Printf("Initial: %s %s", r.URL.Path, bytes.ReplaceAll(reqBody, []byte("\n"), []byte(" ")))

	// Call hook in tradingview using fields available in tv-alerts, does have duplicate info but is based off provided fields
	// https://server/RP/exchange:pair/base/quote/timeframe

	// Create a new alert from the url and request body
	// TODO make this a constructor function
	alert := Alert{
		Exchange:  strings.Split(strings.Split(r.URL.Path, "/")[2], ":")[0],
		Base:      strings.Split(r.URL.Path, "/")[3],
		Quote:     strings.Split(r.URL.Path, "/")[4],
		Timeframe: strings.Split(r.URL.Path, "/")[5],
	}

	// Organize indicator data before converting to json
	// Seperate first line and rest of message and convert to lowercase string
	signal, body, _ := strings.Cut(strings.ToLower(string(reqBody)), "\n")

	// If signal contains MoneyBag, set Call to takeprofit
	const MoneyBag = "💰"

	// map body to usable data json
	body = "{" + body + "}"
	var bodyMap map[string]interface{}
	json.Unmarshal([]byte(body), &bodyMap)

	if strings.Contains(signal, MoneyBag) {

		// this is all broken but I am going to redo it not using signal after the take profit signal since the data source is inconsistent
		// set alert Call to takeprofit
		alert.Call = "EXIT"
		alert.Type = "TAKE"

		// split the body into lines
		pnl, bars, _ := strings.Cut(body, "\n")
		// remove non numbers from parts and cast to float64/int
		alert.Bars, _ = strconv.Atoi(KeepNumbersReg.ReplaceAllString(bars, ""))
		// set price from the signal message
		alert.Price, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(signal, ":")[1]), 64)
		// set pnl from the body
		alert.PNL, _ = strconv.ParseFloat(KeepNumbersReg.ReplaceAllString(pnl, ""), 64)

	} else if strings.Contains(signal, "entry") {

		// set Call to entry
		alert.Call = "ENTRY"

		// set Type to long or short
		if strings.Contains(body, "long") {
			alert.Type = "LONG"
		} else if strings.Contains(body, "short") {
			alert.Type = "SHORT"
		} else {
			logit.Log.Println("Unknown signal type")
		}

		// set price from body map
		alert.Price, _ = strconv.ParseFloat(bodyMap["enter"].(string), 64)
		alert.StopLoss, _ = strconv.ParseFloat(bodyMap["sl"].(string), 64)
		alert.StopLossPercent, _ = strconv.ParseFloat(bodyMap["slp"].(string), 64)

		// fmt.Println("bodyMap: ", bodyMap)

	} else if strings.Contains(signal, "exit") {
		// do other exit Call stuff
		alert.Call = "EXIT"

		// set Type to long or short
		if strings.Contains(body, "long") {
			alert.Type = "LONG"
		} else if strings.Contains(body, "short") {
			alert.Type = "SHORT"
		} else if strings.Contains(body, "stop") {
			alert.Type = "STOP"
		} else {
			logit.Log.Println("Unknown signal type")
		}

	} else if strings.Contains(signal, "error") {
		// If signal contains error, set Call to error
		alert.Call = "ERROR"
	}

	// Convert alert to json
	alertBytes, _ := json.Marshal(alert)

	w.Write([]byte(alertBytes))

	//fmt.Printf("%#v\n", body)
	//fmt.Printf("%#v\n", signal)

	fmt.Printf("%+v\n", alert)
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	// this can be used to process indicators that allow you to set the message data as json
	// right now it doesn't do anything but play with the data

	//fmt.Printf("Headers: %v\n", r.Header)
	webhookData := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&webhookData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Printf("%#v\n", webhookData)
	fmt.Println("webhookData: ", webhookData)

	// for k, v := range webhookData {
	// 	fmt.Println("key: ", k, "value: ", v)
	// }

	webhookDataBytes, _ := json.Marshal(webhookData)
	w.Write([]byte(webhookDataBytes))
}

func main() {

	// Create a new server
	// TODO Use https://github.com/julienschmidt/httprouter for routing

	logit.Log.Println("Listening on port 8080")
	logit.Log.Printf("Server v%s pid=%d started with processes: %d", VERSION, os.Getpid(), runtime.GOMAXPROCS(runtime.NumCPU()))

	http.HandleFunc("/RP/", handleRP)
	http.HandleFunc("/json/", handleJSON)
	logit.Log.Fatal(http.ListenAndServe(":8080", nil))
}
