package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"fmt"
	"os"
	"log"
	"path/filepath"
)

var RootCmd = &cobra.Command{
	Use: "zmeyCoin",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ZmeyCoin client")
		cmd.Usage()
	},
}
var StartWebClient = &cobra.Command{
	Use: "start",
	Short: "Start the zmeyCoin node",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting zmey client")
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.AddCommand(StartWebClient)
}

func initConfig() {
	// Find home directory.
	exePath, err  := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	home := filepath.Dir(exePath)
	fmt.Println(home)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
