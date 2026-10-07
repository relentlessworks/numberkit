package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr   string
	NoAuth bool
}

func Load() *Config {
	c := &Config{
		Addr:   ":7780",
		NoAuth: true,
	}

	// Env vars
	if v := os.Getenv("NUMBERKIT_ADDR"); v != "" {
		c.Addr = v
	}
	if v := os.Getenv("NUMBERKIT_NO_AUTH"); v == "true" || v == "1" {
		c.NoAuth = true
	}

	// Flags
	flag.StringVar(&c.Addr, "addr", c.Addr, "listen address")
	flag.BoolVar(&c.NoAuth, "no-auth", c.NoAuth, "disable authentication")
	flag.Parse()

	return c
}
