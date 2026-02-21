package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "nogocli",
		Usage:       "NogoChain command line interface",
		Description: "Interact with NogoChain nodes and manage accounts",
		Version:     "1.0.0",
		Commands: []*cli.Command{
			{
				Name:    "block",
				Aliases: []string{"b"},
				Usage:   "Block related commands",
				Subcommands: []*cli.Command{
					{
						Name:    "get",
						Aliases: []string{"g"},
						Usage:   "Get block by number or hash",
						Action: func(c *cli.Context) error {
							hashOrNumber := c.Args().First()
							if hashOrNumber == "" {
								return fmt.Errorf("block hash or number is required")
							}
							fmt.Printf("Getting block: %s\n", hashOrNumber)
							return nil
						},
					},
					{
						Name:    "latest",
						Aliases: []string{"l"},
						Usage:   "Get latest block",
						Action: func(c *cli.Context) error {
							fmt.Println("Getting latest block...")
							return nil
						},
					},
					{
						Name:    "height",
						Aliases: []string{"h"},
						Usage:   "Get current block height",
						Action: func(c *cli.Context) error {
							fmt.Println("Getting current block height...")
							return nil
						},
					},
				},
			},
			{
				Name:    "tx",
				Aliases: []string{"t"},
				Usage:   "Transaction related commands",
				Subcommands: []*cli.Command{
					{
						Name:    "get",
						Aliases: []string{"g"},
						Usage:   "Get transaction by hash",
						Action: func(c *cli.Context) error {
							txHash := c.Args().First()
							if txHash == "" {
								return fmt.Errorf("transaction hash is required")
							}
							fmt.Printf("Getting transaction: %s\n", txHash)
							return nil
						},
					},
					{
						Name:    "send",
						Aliases: []string{"s"},
						Usage:   "Send transaction",
						Action: func(c *cli.Context) error {
							to := c.Args().Get(0)
							amount := c.Args().Get(1)
							if to == "" || amount == "" {
								return fmt.Errorf("recipient address and amount are required")
							}
							fmt.Printf("Sending %s NOGO to %s\n", amount, to)
							return nil
						},
					},
				},
			},
			{
				Name:    "account",
				Aliases: []string{"a"},
				Usage:   "Account related commands",
				Subcommands: []*cli.Command{
					{
						Name:    "balance",
						Aliases: []string{"b"},
						Usage:   "Get account balance",
						Action: func(c *cli.Context) error {
							address := c.Args().First()
							if address == "" {
								return fmt.Errorf("account address is required")
							}
							fmt.Printf("Getting balance for: %s\n", address)
							return nil
						},
					},
					{
						Name:    "new",
						Aliases: []string{"n"},
						Usage:   "Create new account",
						Action: func(c *cli.Context) error {
							fmt.Println("Creating new account...")
							return nil
						},
					},
				},
			},
			{
				Name:    "network",
				Aliases: []string{"n"},
				Usage:   "Network related commands",
				Subcommands: []*cli.Command{
					{
						Name:    "status",
						Aliases: []string{"s"},
						Usage:   "Get network status",
						Action: func(c *cli.Context) error {
							fmt.Println("Getting network status...")
							return nil
						},
					},
					{
						Name:    "peers",
						Aliases: []string{"p"},
						Usage:   "Get peer list",
						Action: func(c *cli.Context) error {
							fmt.Println("Getting peer list...")
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
