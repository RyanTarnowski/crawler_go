package main

import (
	"fmt"
	"os"
)

func main() {
	//Get cmd line args, first arg is the program path so start at the second index
	//test URL: https://learnwebscraping.dev/practice/ecommerce/
	//test URL: https://learnwebscraping.dev/practice/ecommerce/products/ashenfang-longsword-fan-1001/
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawBaseURL := args[0]

	fmt.Printf("starting crawl of: %v\n", rawBaseURL)

	rawHTML, err := getHTML(rawBaseURL)
	if err != nil {
		fmt.Printf("Failed to get HTML: %v\n", err)
	}

	fmt.Println(rawHTML)
}
