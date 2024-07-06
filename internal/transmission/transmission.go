package transmission

import (
	"log"
	"strconv"

	"github.com/hekmon/transmissionrpc"
	"github.com/spf13/viper"
)

func Main() {
	host := viper.GetString("transmission.host")
	portStr := viper.GetString("transmission.port")
	user := viper.GetString("transmission.user")
	password := viper.GetString("transmission.password")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("[ERROR] Invalid port number: %v", err)
	}

	// Connecting to transmission
	client, err := transmissionrpc.New(host, user, password, &transmissionrpc.AdvancedConfig{Port: uint16(port)})
	if err != nil {
		log.Fatalf("[ERROR] Error creating transmission client: %v", err)
	}

	stats, err := client.SessionStats()
	if err != nil {
		log.Fatalf("[ERROR] Error getting session stats: %v", err)
	}

	log.Println(stats)
}
