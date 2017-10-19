package main

import (
	"fmt"
	"github.com/client9/go-disconnectme"
)

func main() {
	dm, err := disconnectme.ParseURL()
	if err != nil {
		panic(err)
	}
	for category, vendors := range dm {
		fmt.Printf("Category %q has %d entries\n", category, len(vendors))
		for _, vendor := range vendors {
			fmt.Printf("  %s --> %s\n", vendor.Name, vendor.Address)
			for _, domain := range vendor.Domains {
				fmt.Printf("      %s\n", domain)
			}
		}
	}
}
