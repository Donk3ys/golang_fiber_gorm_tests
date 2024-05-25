package storage

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/valkey-io/valkey-go"
)

func TestConnectValkey(ctx context.Context) (valkey.Client, testcontainers.Container) {
	// host := os.Getenv("VALKEY_HOST")
	// port := os.Getenv("VALKEY_PORT")
	// client, err := valkey.NewClient(valkey.ClientOption{
	// 	InitAddress: []string{host + ":" + port},
	// 	Username:    os.Getenv("VALKEY_USER"),
	// 	Password:    os.Getenv("VALKEY_PASSWORD"),
	// })
	// if err != nil {
	// 	log.Fatalf("ValKey connection error %s:%s\n", host, port, err)
	// 	return nil
	// }

	req := testcontainers.ContainerRequest{
		Image:        "valkey/valkey:alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	tc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		// Reuse:            true,
	})
	if err != nil {
		log.Fatalf("Could not start Valkey: %s", err)
	}

	addr, err := tc.Endpoint(ctx, "")
	if err != nil {
		log.Fatalf("Unable to retrieve the endpoint: %s", err)
	}

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{addr},
	})
	if err != nil {
		log.Panic("ValKey test container connection error\n", err)
	}

	return client, tc
}
