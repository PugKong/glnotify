package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type GitlabFake struct {
	mergeRequests      []MergeRequest
	mergeRequestsNotes map[string][]Note
}

func NewGitlabFake() *GitlabFake {
	return &GitlabFake{
		mergeRequestsNotes: make(map[string][]Note),
	}
}

func (g *GitlabFake) AddMergeRequest(mergeRequest MergeRequest) {
	g.mergeRequests = append(g.mergeRequests, mergeRequest)
}

func (g *GitlabFake) AddMergeRequestNote(mergeRequest MergeRequest, note Note) {
	key := createMergeRequestKey(mergeRequest)
	g.mergeRequestsNotes[key] = append(g.mergeRequestsNotes[key], note)
}

func (g *GitlabFake) ListOpenedMergeRequests() ([]MergeRequest, error) {
	return g.mergeRequests, nil
}

func (g *GitlabFake) ListMergeRequestNotes(mergeRequest MergeRequest) ([]Note, error) {
	key := createMergeRequestKey(mergeRequest)

	return g.mergeRequestsNotes[key], nil
}

func TestApp_Run(t *testing.T) {
	t.Run("it prints information about new comments", func(t *testing.T) {
		mergeRequest := MergeRequest{
			ProjectID: 1,
			IID:       2,
			URL:       "https://example.com",
			Title:     "MR Title",
		}

		gitlab := NewGitlabFake()
		gitlab.AddMergeRequest(mergeRequest)
		gitlab.AddMergeRequestNote(mergeRequest, Note{ID: 3, Author: "John Doe"})
		gitlab.AddMergeRequestNote(mergeRequest, Note{ID: 4, Author: "Jane Doe"})
		gitlab.AddMergeRequestNote(mergeRequest, Note{ID: 5, Author: "John Doe"})

		state := NewState()
		out := bytes.NewBuffer([]byte{})
		want := fmt.Sprintf("MR: %s\nCommented by: %s\n%s\n\n", mergeRequest.Title, "Jane Doe, John Doe", mergeRequest.URL)

		app := NewApp(out, gitlab)
		state, err := app.Run(state)
		require.NoError(t, err)
		require.Equal(t, want, out.String())

		out.Reset()
		want = ""

		_, err = app.Run(state)
		require.NoError(t, err)
		require.Equal(t, want, out.String())
	})

	t.Run("it prints information about new labels", func(t *testing.T) {
		mergeRequest := MergeRequest{
			ProjectID: 1,
			IID:       2,
			URL:       "https://example.com",
			Title:     "MR Title",
			Labels:    []string{"foo", "bar"},
		}

		gitlab := NewGitlabFake()
		gitlab.AddMergeRequest(mergeRequest)

		state := NewState()
		out := bytes.NewBuffer([]byte{})
		want := fmt.Sprintf(
			"MR: %s\nLabeled: %s\n%s\n\n",
			mergeRequest.Title,
			strings.Join(mergeRequest.Labels, ", "),
			mergeRequest.URL,
		)

		app := NewApp(out, gitlab)
		state, err := app.Run(state)
		require.NoError(t, err)
		require.Equal(t, want, out.String())

		out.Reset()
		want = ""

		_, err = app.Run(state)
		require.NoError(t, err)
		require.Equal(t, want, out.String())
	})
}
