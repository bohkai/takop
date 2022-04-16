package main

import (
	_ "embed"
	"encoding/json"
	"errors"
)

type Config struct {
	Google *GoogleConfig `json:"google"`
	Discord *DiscordConfig `json:"discord"`
}

type GoogleConfig struct {
	Key string `json:"key"`
	ID string `json:"id"`
}

type DiscordConfig struct {
	Token string `json:"token"`
}

//go:embed config.json
var configBytes []byte

func NewConfig() (*Config, error) {
	if len(configBytes) < 1 {
		return nil, errors.New("config not loaded")
	}

	var c *Config
	err := json.Unmarshal(configBytes, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}