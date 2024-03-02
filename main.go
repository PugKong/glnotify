package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/xanzy/go-gitlab"
	"golang.org/x/time/rate"
)

const (
	rateLimitR = rate.Limit(2)
	rateLimitB = 1
)

type Config struct {
	BaseURL    string `json:"base_url"`
	Token      string `json:"token"`
	UserID     int    `json:"user_id"`
	ProjectIDs []int  `json:"project_ids"`
}

func main() {
	var _ Gitlab = (*XanzyGitlab)(nil)

	homePath, err := os.UserConfigDir()
	must(err)

	config, err := loadConfig(path.Join(homePath, "glnotify", "config.json"))
	must(err)

	statePath := path.Join(homePath, "glnotify", "state.json")
	state, err := loadState(statePath)
	must(err)

	client, err := gitlab.NewClient(
		config.Token,
		gitlab.WithBaseURL(config.BaseURL),
		gitlab.WithCustomLimiter(rate.NewLimiter(rateLimitR, rateLimitB)),
	)
	must(err)

	gitlab := NewXanzyGitlab(client, config.UserID, config.ProjectIDs)

	app := NewApp(gitlab)
	state, err = app.Run(state, os.Stdout)
	must(err)

	must(saveState(statePath, state))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func loadConfig(path string) (Config, error) {
	var config Config

	file, err := os.Open(path)
	if err != nil {
		return config, fmt.Errorf("unable to open %q config file: %w", path, err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return config, fmt.Errorf("unable to parse %q config file: %w", path, err)
	}

	return config, nil
}

func loadState(path string) (State, error) {
	state := NewState()

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return state, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return state, fmt.Errorf("unable to open %q state file: %w", path, err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&state)
	if err != nil {
		return state, fmt.Errorf("unable to parse %q state file: %w", path, err)
	}

	return state, nil
}

func saveState(path string, state State) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create/truncate %q state file: %w", path, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(state)
	if err != nil {
		return fmt.Errorf("unable to save %q state file: %w", path, err)
	}

	return nil
}
