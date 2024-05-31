package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	Port        int
	Backends    []string
	STATIC_ROOT string
	MEDIA_ROOT  string
	ALGO        string
}

var GLOBAL Config

func ConfigLog() {
	Handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	Logger = slog.New(Handler)
}

func InitialiseSystem() error {
	ConfigLog()
	Logger.Info("Configuring the Go-Balancer")
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&GLOBAL)
	if err != nil {
		return err
	}
	Logger.Info("Configuration Successfully Parsed")
	n := len(GLOBAL.Backends)
	Logger.Info(fmt.Sprintf("Available Backends to balance: %d", n))

	if !check_roots() {
		Logger.Error("Error in reading the STATIC and MEDIA roots")
		return fmt.Errorf("either STATIC or MEDIA or BOTH roots are Invalid")
	}

	switch GLOBAL.ALGO {
	case "RR":
		MANAGER = GetPool(0)
	case "LC":
		MANAGER = GetQueue(0)
	}
	Logger.Info(fmt.Sprintf("Chosen Scheduling Algorithm: %s", GLOBAL.ALGO))

	for _, url := range GLOBAL.Backends {
		backend := GetBackend(url)
		MANAGER.Addserver(backend)
	}
	return nil
}

func check_roots() bool {
	if _, err := os.Stat(GLOBAL.STATIC_ROOT); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(GLOBAL.MEDIA_ROOT); os.IsNotExist(err) {
		return false
	}
	return true
}
