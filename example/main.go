package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func logJSON(v interface{}) {
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("Error marshaling to JSON: %v", err)
		return
	}
	log.Printf("%s", string(jsonBytes))
}

func printUsage() {
	fmt.Println("ZeroDev Go SDK Examples")
	fmt.Println("\nUsage:")
	fmt.Println("  go run . <command>")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  4337    Run 4337 account example")
	fmt.Println("  7702    Run 7702 account example (full flow with authorization)")
	fmt.Println("  help       Show this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  go run . 7702")
	fmt.Println("\nMake sure to set up your .env file with ZERODEV_PROJECT_ID and USEROP_BUILDER_API_KEY")
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Error: No command specified")
		printUsage()
		os.Exit(1)
	}

	command := args[0]

	switch command {
	case "7702":
		run4337Example()
	case "4337":
		run7702Example()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Error: Unknown command '%s'\n\n", command)
		printUsage()
		os.Exit(1)
	}
}
