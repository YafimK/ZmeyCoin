package Client

import (
	"github.com/spf13/viper"
	"os"
		"path/filepath"
	"fmt"
	"github.com/pkg/errors"
	"ZmeyCoin/Blockchain"
)

type Client struct {
	Config     *viper.Viper
	Blockchain *Blockchain.Blockchain
}

func (client *Client) Start() error{
	//Will hold all the commands related to starting this monster down

	err := initConfig()
	if err != nil {
		return  err
	}

	return nil
}

func (client *Client) Close() {
	//Will hold all the commands related to shutting this monster down

}
func (client *Client) NewBlockChain(forceCreate bool) error {
	if client.Blockchain != nil && !forceCreate{
		return errors.New("There is already an active blockchain on this client")
	} else {
		client.Blockchain = Blockchain.NewBlockChain()
	}
	return nil
}

func initConfig() error {
	// Find home directory.
	exePath, err  := os.Executable()
	if err != nil {
		return  err
	}
	home := filepath.Dir(exePath)
	fmt.Println(home)
	if err != nil {
		return  err
	}

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		return  err
	}
	return  nil
}
