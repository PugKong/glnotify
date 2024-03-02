package main

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/xanzy/go-gitlab"
	"golang.org/x/time/rate"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type Config struct {
	BaseUrl    string `json:"base_url"`
	Token      string `json:"token"`
	UserId     int    `json:"user_id"`
	ProjectIds []int  `json:"project_ids"`
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
		gitlab.WithBaseURL(config.BaseUrl),
		gitlab.WithCustomLimiter(rate.NewLimiter(rate.Limit(2), 1)),
	)
	must(err)

	gitlab := NewXanzyGitlab(client, config.UserId, config.ProjectIds)

	app := NewApp(gitlab)
	state, err = app.Run(state)
	must(err)

	must(saveState(statePath, state))
}

func loadConfig(path string) (Config, error) {
	var config Config

	f, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&config)
	return config, err
}

func loadState(path string) (State, error) {
	state := NewState()

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return state, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return state, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&state)
	return state, err
}

func saveState(path string, state State) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	return encoder.Encode(state)
}
