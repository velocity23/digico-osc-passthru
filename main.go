package main

import (
	"digico-osc-passthru/osc"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Starting DiGiCo OSC Passthru...")

	clients := []osc.Client{}

	consoleIp, consoleRx, consoleTx := os.Args[1], os.Args[2], os.Args[3]

	txInt, err := strconv.Atoi(consoleTx)
	if err != nil {
		fmt.Println("Error converting tx port to int")
		return
	}
	consoleClient := osc.NewClient(consoleIp, txInt)

	d1 := osc.NewStandardDispatcher()
	d1.AddMsgHandler("*", func(msg *osc.Message, addr *net.Addr) {
		consoleClient.Send(msg)
		clients = append(clients, *osc.NewClient((*addr).String(), txInt))
	})
	distroServer := &osc.Server{
		Addr:       ":" + consoleRx,
		Dispatcher: d1,
	}
	distroServer.ListenAndServe()
	defer distroServer.CloseConnection()

	d2 := osc.NewStandardDispatcher()
	d2.AddMsgHandler("*", func(msg *osc.Message, _ *net.Addr) {
		for _, client := range clients {
			client.Send(msg)
		}
	})
	consoleServer := &osc.Server{
		Addr:       ":" + consoleTx,
		Dispatcher: d2,
	}
	consoleServer.ListenAndServe()
	defer consoleServer.CloseConnection()
}
