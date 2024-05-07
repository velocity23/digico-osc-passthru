package main

import (
	"digico-osc-passthru/osc"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Starting DiGiCo OSC Passthru...")

	consoleIp, consoleRx, consoleTx := os.Args[1], os.Args[2], os.Args[3]
	clients := []osc.Client{}

	fmt.Println("RX", consoleRx)
	fmt.Println("TX", consoleTx)

	txInt, err := strconv.Atoi(consoleTx)
	if err != nil {
		fmt.Println("Error converting tx port to int")
		return
	}

	rxInt, err := strconv.Atoi(consoleRx)
	if err != nil {
		fmt.Println("Error converting rx port to int")
		return
	}

	go initConsole(txInt, &clients)
	initClients(consoleIp, rxInt, txInt, &clients)
}

func initConsole(txInt int, clients *[]osc.Client) {
	d := osc.NewStandardDispatcher()
	// From Console
	d.AddMsgHandler("*", func(msg *osc.Message, _ *net.Addr) {
		fmt.Println("From Console", msg.Address)
		if len(*clients) == 0 {
			fmt.Println("No Clients")
		}
		for _, client := range *clients {
			fmt.Println("Sending to", client.IP())
			client.Send(msg)
		}
		fmt.Println("")
	})
	fromConsole := &osc.Server{
		Addr:       "0.0.0.0:" + fmt.Sprint(txInt),
		Dispatcher: d,
	}
	fromConsole.ListenAndServe()
	defer fromConsole.CloseConnection()
}

func initClients(consoleIp string, rxInt int, txInt int, clients *[]osc.Client) {
	toConsole := osc.NewClient(consoleIp, rxInt)
	d := osc.NewStandardDispatcher()
	// From Clients
	d.AddMsgHandler("*", func(msg *osc.Message, addr *net.Addr) {
		toConsole.Send(msg)
		
		for _, client := range *clients {
			if (strings.Split(client.IP(), ":")[0] == strings.Split((*addr).String(), ":")[0]) {
				return
			}
		}
		*clients = append(*clients, *osc.NewClient(strings.Split((*addr).String(), ":")[0], txInt))
	})
	fromClients := &osc.Server{
		Addr:       ":" + fmt.Sprint(rxInt),
		Dispatcher: d,
	}
	fromClients.ListenAndServe()
	defer fromClients.CloseConnection()
}