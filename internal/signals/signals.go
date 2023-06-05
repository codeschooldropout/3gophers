package signals

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codeschooldropout/3gophers/internal/logit"
	"github.com/gen2brain/beeep"
)

func init() {
	// fmt.Println("cmd/signals.go")
}

type Asset struct {
	Exchange  string `json:"exchange"`  // 1. kucoin 2. kraken 3. coinbase - these are the ones i want to start with if possible.
	Base      string `json:"base"`      // BTC, ETH, etc
	Quote     string `json:"quote"`     // USD, EUR, etc
	Timeframe string `json:"timeframe"` // 1m, 5m, 1h, 1d, 1w, 1M, 1y
}
type Alert struct {
	Order     string  `json:"order"`               // entry, exit, error
	Contracts float64 `json:"contracts,omitempty"` // long, short, take
	Price     float64 `json:"price,omitempty"`     // buy high, sell low right?
	Ticker    string  `json:"ticker,omitempty"`    // profit/loss from exit signals
	Interval  int     `json:"interval,omitempty"`  // Number of bars from entry signal to exit signal (bars x timeframe to identify corresponding signals)
	Position  float64 `json:"position,omitempty"`  // Stop loss price to set on entry signal
	TimeNow   string  `json:"timenow,omitempty"`   // Stop loss percent to set on entry signal
	Asset     Asset   `json:"asset,omitempty"`     // Asset to trade
}

func HandleTradingViewJSON(w http.ResponseWriter, r *http.Request) {
	// Process the signal from TradingView

	// Create a new alert to handle incoming data
	var alert Alert

	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// print the alert to the console
	fmt.Printf("alert: %v\n", alert)

	// Convert the alert back to JSON for the response
	alertBytes, err := json.Marshal(alert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// respond with 200 OK
	w.WriteHeader(http.StatusOK)

	// w.Write([]byte(alertBytes))

	//fmt.Printf("Headers: %v\n", r.Header)
	// webhookData := make(map[string]interface{})

	// Output Examples
	//fmt.Printf("%#v\n", webhookData)
	// fmt.Println("webhookData: ", webhookData)

	// for k, v := range webhookData {
	// 	fmt.Println("key: ", k, "value: ", v)
	// }

	// webhookDataBytes, _ := json.Marshal(webhookData)
	// w.Write([]byte(webhookDataBytes))

	err2 := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
	if err2 != nil {
		panic(err)
	}

	logit.Log.Printf("alertBytes: %s", alertBytes)
}

// create a new alert
func NewAlert(order string, contracts float64, price float64, ticker string, interval int, position float64, timenow string) *Alert {
	return &Alert{
		Order:     order,
		Contracts: contracts,
		Price:     price,
		Ticker:    ticker,
		Interval:  interval,
		Position:  position,
		TimeNow:   timenow,
		// Asset:     asset,
	}
}

// // create a new asset
// func NewAsset(exchange string, base string, quote string, timeframe string) *Asset {
// 	return &Asset{
// 		Exchange:  exchange,
// 		Base:      base,
// 		Quote:     quote,
// 		Timeframe: timeframe,
// 	}
// }
