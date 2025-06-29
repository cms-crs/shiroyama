package config

import "time"

type ClientsConfig struct {
	BoardServiceAddr string
	TeamServiceAddr  string
	DialTimeout      time.Duration
}
