package main

import (
	_ "embed"
	"encoding/json"
	"errors"
)

type DiscordConfig struct {
	Token string `json:"token"`
}

//go:embed config.json
var configBytes []byte

func NewConfig() (*DiscordConfig, error) {
	if len(configBytes) < 1 {
		return nil, errors.New("config not loaded")
	}

	var c *DiscordConfig
	err := json.Unmarshal(configBytes, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}