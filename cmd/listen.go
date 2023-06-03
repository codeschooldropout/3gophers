/*
Copyright © 2022 codeschooldropout code@cay.io
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/codeschooldropout/3gophers/internal/logit"
	"github.com/codeschooldropout/3gophers/internal/signals"
	"github.com/codeschooldropout/3gophers/internal/ui"
	"github.com/spf13/cobra"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Starts the listening server to capture calls from tradingview",
	Long:  `Starts the listening server to capture indicator calls from tradingview and process their signals `,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("args %v", args)
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		logit.Log.Println("listen called on port", port)

		go startHTTPServer(port)

		p := tea.NewProgram(ui.NewModel())

		if _, err := p.Run(); err != nil {
			logit.Log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
	listenCmd.PersistentFlags().IntP("port", "p", 8080, "Run listener on specified port")
}

func startHTTPServer(port int) {

	// Create a new server
	logit.Log.Println("Listening on port", port)
	logit.Log.Printf("Server v%s pid=%d started with processes: %d", Version, os.Getpid(), runtime.GOMAXPROCS(runtime.NumCPU()))

	// http.HandleFunc("/RP/", signals.HandleTradingViewRP)
	http.HandleFunc("/json/", signals.HandleTradingViewJSON)
	logit.Log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
