package main

import (
	"pkg/cli"
	"pkg/db"
)

func main() {
	defer db.Close()
	cli.Start()
	// wallet.Wallet()
}
