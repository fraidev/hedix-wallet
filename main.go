package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fraidev/hedix-wallet/models"
	"github.com/fraidev/hedix-wallet/services"
)

func main() {
	// Create wallet
	wallet := services.NewWallet()

	// Check if user wants to use file or interactive mode (default)
	if len(os.Args) > 1 && os.Args[1] == "--file" {
		pathName := os.Args[2]
		file, err := os.Open(pathName)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}
		defer file.Close()

		runFile(wallet, file)
	} else {
		runInteractive(wallet)
	}
}

func runFile(wallet *services.Wallet, file *os.File) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input := scanner.Text()
		tx, err := models.ParseTransaction(input)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		err = wallet.ProcessTransaction(tx)
		if err != nil {
			fmt.Printf("Transaction failed: %s\n", err)
		}

		fmt.Printf("   State: %s\n", wallet)
	}

	fmt.Println()
	fmt.Printf("Final Balance: %s\n", wallet)
}

func runInteractive(wallet *services.Wallet) {
	fmt.Println("Interactive Mode - Enter transactions")
	fmt.Println("Format: <DEPOSIT|WITHDRAW> <BTC|ETH|USD> <amount>")
	fmt.Println("Example: DEPOSIT BTC 1.5")
	fmt.Println()
	fmt.Printf("Current State: %s\n", wallet)
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		tx, err := models.ParseTransaction(input)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		err = wallet.ProcessTransaction(tx)
		if err != nil {
			fmt.Printf("Transaction failed: %s\n", err)
		} else {
			fmt.Println("Transaction successful")
		}

		fmt.Printf("Current State: %s\n", wallet)
	}

	fmt.Println()
	fmt.Printf("Final Balance: %s\n", wallet)
}
