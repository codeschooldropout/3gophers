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
	VERSION = "0.0.2"
)

type Asset struct {
	Exchange  string `json:"exchange"`  // 1. kucoin 2. kraken 3. coinbase - these are the ones i want to start with if possible.
	Base      string `json:"base"`      // BTC, ETH, etc
	Quote     string `json:"quote"`     // USD, EUR, etc
	Timeframe string `json:"timeframe"` // 1m, 5m, 1h, 1d, 1w, 1M, 1y
}
type Alert struct {
	Call            string  `json:"call"`               // entry, exit, error
	Position        string  `json:"position,omitempty"` // long, short, take
	Price           float64 `json:"price,omitempty"`    // buy high, sell low right?
	PNL             float64 `json:"pnl,omitempty"`      // profit/loss from exit signals
	Bars            int     `json:"bars,omitempty"`     // Number of bars from entry signal to exit signal (bars x timeframe to identify corresponding signals)
	StopLoss        float64 `json:"sl,omitempty"`       // Stop loss price to set on entry signal
	StopLossPercent float64 `json:"slp,omitempty"`      // Stop loss percent to set on entry signal
	Asset           Asset   `json:"asset,omitempty"`    // Asset to trade
}

func handleRP(w http.ResponseWriter, r *http.Request) {
	// create regex to match numbers and certain characters
	keepNumbersReg, err := regexp.Compile("[^0-9-.]")

	// Symbols to use
	// callIcon := "\u260E" // these icons work too (might look better as well)
	moneyIcon := "💰"
	stoplossIcon := "🛑"
	enterIcon := "✅"
	exitIcon := "❌"
	errorIcon := "🔥"
	longIcon := "📈"
	shortIcon := "📉"

	callIcon, positionIcon := "", ""

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
	asset := Asset{
		Exchange:  strings.Split(strings.Split(r.URL.Path, "/")[2], ":")[0],
		Base:      strings.Split(r.URL.Path, "/")[3],
		Quote:     strings.Split(r.URL.Path, "/")[4],
		Timeframe: strings.Split(r.URL.Path, "/")[5],
	}

	alert := Alert{}
	alert.Asset = asset

	// Seperate first line and rest of message and convert to lowercase string
	signal, body, _ := strings.Cut(strings.ToLower(string(reqBody)), "\n")

	// body is almost well formed for json but needs external brackets
	body = "{" + body + "}"
	// map body so i'm not just strings.containing everything?
	var bodyMap map[string]interface{}
	json.Unmarshal([]byte(body), &bodyMap)

	// Occassionally signals will register from tradingview as a signal error
	// Identify this and return an error
	if signal == "alert error" {
		alert.Call = "error"
		callIcon = errorIcon
		positionIcon = errorIcon
	}

	// EXIT SIGNALS
	// TAKEPROFIT
	// The takeprofit signal does not use the same format, or a well defined format like the other alert message bodies
	// So we need to handle it separately and identify if signal is take profit
	if strings.Contains(signal, moneyIcon) {
		// split the body into lines for parsing
		pnl, bars, _ := strings.Cut(body, "\n")
		// set call and position
		alert.Call = "exit"
		alert.Position = "take"
		// remove non numbers from parts and cast to correct number type
		alert.Bars, _ = strconv.Atoi(keepNumbersReg.ReplaceAllString(bars, ""))
		// set price from the signal message
		alert.Price, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(signal, ":")[1]), 64)
		// set pnl from the body
		alert.PNL, _ = strconv.ParseFloat(keepNumbersReg.ReplaceAllString(pnl, ""), 64)
		// set alert icon
		alert.StopLoss, alert.StopLossPercent = 0, 0
		callIcon = exitIcon
		positionIcon = moneyIcon

	}

	if _, ok := bodyMap["exit"]; ok {
		// EXIT LONG & SHORT & STOPLOSS
		// set call
		alert.Call = "exit"
		// set long/short position
		if bodyMap["type"].(string) == "hard long exit" {
			alert.Position = "long"
			positionIcon = longIcon
		} else if bodyMap["type"].(string) == "hard short exit" {
			alert.Position = "short"
			positionIcon = shortIcon
		} else if bodyMap["type"].(string) == "stop loss hit" {
			alert.Position = "stop"
			// change the icon if its a stop loss
			positionIcon = stoplossIcon
		}

		// set price from the body
		alert.Price, _ = strconv.ParseFloat(bodyMap["exit"].(string), 64)
		alert.PNL, _ = strconv.ParseFloat(keepNumbersReg.ReplaceAllString(bodyMap["pnl"].(string), ""), 64)
		alert.Bars, _ = strconv.Atoi(bodyMap["traded bars"].(string))
		// set the icon
		callIcon = exitIcon

	} else if _, ok := bodyMap["enter"]; ok {
		// ENTER LONG & SHORT
		// set call
		alert.Call = "enter"
		// set long/short position
		if bodyMap["type"].(string) == "long signal" {
			alert.Position = "long"
			positionIcon = longIcon
		} else if bodyMap["type"].(string) == "short signal" {
			alert.Position = "short"
			positionIcon = shortIcon
		}

		// set price from body map
		alert.Price, _ = strconv.ParseFloat(bodyMap["enter"].(string), 64)
		alert.StopLoss, _ = strconv.ParseFloat(bodyMap["sl"].(string), 64)
		alert.StopLossPercent, _ = strconv.ParseFloat(bodyMap["slp"].(string), 64)
		//set the icon
		callIcon = enterIcon
	}

	// Convert alert to json
	alertBytes, _ := json.Marshal(alert)
	// return json to caller (this isn't needed except for testing or if i send signals from things that are not tradingview
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(alertBytes))

	// print alert to console
	logit.Log.Printf("%s %s Alert: %s", callIcon, positionIcon, string(alertBytes))
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
