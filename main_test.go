package main

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Run("it loads config", func(t *testing.T) {
		path := path.Join(t.TempDir(), "config.json")

		file, err := os.Create(path)
		require.NoError(t, err, "failed to create temp file")
		defer file.Close()

		config := Config{
			BaseURL:    "https://example.com",
			Token:      "token",
			UserID:     42,
			ProjectIDs: []int{43, 44},
		}

		err = json.NewEncoder(file).Encode(config)
		require.NoError(t, err, "failed to write to temp file")

		loadedConfig, err := loadConfig(path)
		require.NoError(t, err)
		require.Equal(t, config, loadedConfig)
	})
}

func TestLoadState(t *testing.T) {
	t.Run("it loads state", func(t *testing.T) {
		path := path.Join(t.TempDir(), "state.json")

		file, err := os.Create(path)
		require.NoError(t, err, "failed to create temp file")
		defer file.Close()

		state := NewState()
		state.AddNote(MergeRequest{}, Note{})

		err = json.NewEncoder(file).Encode(state)
		require.NoError(t, err, "failed to write to temp file")

		loadedState, err := loadState(path)
		require.NoError(t, err)
		require.Equal(t, state, loadedState)
	})
}

func TestSaveState(t *testing.T) {
	t.Run("it saves state", func(t *testing.T) {
		path := path.Join(t.TempDir(), "state.json")

		state := NewState()
		state.AddNote(MergeRequest{}, Note{})

		err := saveState(path, state)
		require.NoError(t, err)

		file, err := os.Open(path)
		require.NoError(t, err)
		defer file.Close()

		savedState := NewState()
		err = json.NewDecoder(file).Decode(&savedState)
		require.NoError(t, err)
		require.Equal(t, state, savedState)
	})
}
