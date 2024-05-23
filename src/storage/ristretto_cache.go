package storage

import (
	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"
	"github.com/gofiber/fiber/v2/log"
)

func ConnectRistrettoCache() *marshaler.Marshaler {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     100,
		BufferItems: 64,
	})
	if err != nil {
		log.Fatal("Ristretto cache connection error", err)
	}
	ristrettoStore := ristretto_store.NewRistretto(ristrettoCache)
	cacheMan := cache.New[any](ristrettoStore)
	return marshaler.New(cacheMan)
}
