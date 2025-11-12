package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
)

type Config struct {
	Addr       string `json:"addr"`
	Port       int    `json:"port"`
	User       string `json:"user"`
	PrivateKey string `json:"private_key"`
	Password   string `json:"password"`
}

func Init() *Config {
	c := &Config{
		Addr:       os.Getenv("PVE_ADDR"),
		Port:       22,
		User:       "root",
		PrivateKey: "~/.ssh/id_ed25519",
		Password:   "",
	}

	if c.Addr == "" {
		log.Fatal("pve addr is required")
	}

	if val, ok := os.LookupEnv("PVE_PORT"); ok {
		c.Port = cast.ToInt(val)
	}

	if val, ok := os.LookupEnv("PVE_USER"); ok {
		c.User = val
	}

	if val, ok := os.LookupEnv("PVE_PRIVATE_KEY"); ok {
		c.PrivateKey = val
	}

	if strings.HasPrefix(c.PrivateKey, "~") {
		c.PrivateKey = filepath.Join(os.Getenv("HOME"), c.PrivateKey[1:])
	}

	if val, err := filepath.Abs(c.PrivateKey); err == nil {
		c.PrivateKey = val
	}

	if fi, err := os.Stat(c.PrivateKey); err == nil && !fi.IsDir() {
		if val, err := os.ReadFile(c.PrivateKey); err == nil {
			c.PrivateKey = string(val)
		}
	}

	if val, ok := os.LookupEnv("PVE_PRIVATE_KEY_PASSWORD"); ok {
		c.Password = val
	}

	return c
}
