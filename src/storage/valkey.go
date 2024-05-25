package storage

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/valkey-io/valkey-go"
)

func ConnectValkey() valkey.Client {
	host := os.Getenv("VALKEY_HOST")
	port := os.Getenv("VALKEY_PORT")
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{host + ":" + port},
		Username:    os.Getenv("VALKEY_USER"),
		Password:    os.Getenv("VALKEY_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("ValKey connection error %s:%s\n", host, port, err)
		return nil
	}
	log.Infof("ValKey connected %s:%s ", host, port)

	return client
}
