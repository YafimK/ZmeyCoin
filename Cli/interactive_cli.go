package Cli

import (
	"errors"
	"github.com/abiosoft/ishell"
	"ZmeyCoin/Client"
)

func Cli(){
	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()
	client := Client.Client{}
	err := client.Start()
	if err != nil {
		panic("Failed starting zmeyCoin client instance")
	}
	shell.Println("ZmeyCoin interactive cli")

	shell.AddCmd(&ishell.Cmd{
		Name: "init blockchain",
		Help: "Creates a new blockchain with genesis block inside",
		Func: func(c *ishell.Context) {
			err := client.NewBlockChain(false)
			if err != nil {
				c.ShowPrompt(false)
				defer c.ShowPrompt(true)
				c.Println("There is already a blockchain on this client would you like to reset it ? (y/n)")
				answer := c.ReadLine()
				if answer == "y" {
					c.Println("Creating new blockchain in a moment")
					err := client.NewBlockChain(true)
					if err!=nil {
						panic(err)
					}
				}
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "get balance",
		Help: "get balance for specific address, you can press tab for autocomplete",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 0 {
				c.Err(errors.New("you have to provide the adderss for which you'd like the balance for"))
				c.Println(c.HelpText())
			}
		},
		Completer: func([]string) []string {
			return []string{} //TODO: get all addresses from wallet
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "create wallet",
		Help: "Generates a new key-pair and saves it into the wallet file",
		Func: func(c *ishell.Context) {

		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "print blockchain",
		Help: "prints the current blockchain",
		Func: func(c *ishell.Context) {

		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "list addresses",
		Help: "lists all wallet addresses available at this client disposal",
		Func: func(c *ishell.Context) {

		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "send",
		Help: "sends coins from one address to another, i.e. creates a transaction",
		Func: func(c *ishell.Context) {

		},
	})

	// run shell
	shell.Run()
}


func main (){
	Cli()
}