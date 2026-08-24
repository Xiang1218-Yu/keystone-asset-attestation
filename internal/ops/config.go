package ops

import (
	"os"
	"strconv"
)

type Config struct {
	Address string
	Workers int
}

func LoadConfig() Config {
	address := os.Getenv("SERVICE_ADDR")
	if address == "" {
		address = ":8282"
	}
	workers, err := strconv.Atoi(os.Getenv("SERVICE_WORKERS"))
	if err != nil || workers < 1 {
		workers = 3
	}
	return Config{Address: address, Workers: workers}
}
