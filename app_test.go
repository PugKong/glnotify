package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type GitlabFake struct {
	mergeRequests      []MergeRequest
	mergeRequestsNotes map[MergeRequest][]Note
}

func NewGitlabFake() *GitlabFake {
	return &GitlabFake{
		mergeRequestsNotes: make(map[MergeRequest][]Note),
	}
}

func (g *GitlabFake) AddMergeRequest(mergeRequest MergeRequest) {
	g.mergeRequests = append(g.mergeRequests, mergeRequest)
}

func (g *GitlabFake) AddMergeRequestNote(mergeRequest MergeRequest, note Note) {
	g.mergeRequestsNotes[mergeRequest] = append(g.mergeRequestsNotes[mergeRequest], note)
}

func (g *GitlabFake) ListOpenedMergeRequests() ([]MergeRequest, error) {
	return g.mergeRequests, nil
}

func (g *GitlabFake) ListMergeRequestNotes(mergeRequest MergeRequest) ([]Note, error) {
	return g.mergeRequestsNotes[mergeRequest], nil
}

func TestApp_Run(t *testing.T) {
	t.Run("it prints information about new comments", func(t *testing.T) {
		mergeRequest := MergeRequest{
			ProjectID: 1,
			IID:       2,
			URL:       "https://example.com",
			Title:     "MR Title",
		}
		note := Note{
			ID:     3,
			Author: "John Doe",
		}

		gitlab := NewGitlabFake()
		gitlab.AddMergeRequest(mergeRequest)
		gitlab.AddMergeRequestNote(mergeRequest, note)

		state := NewState()
		out := bytes.NewBuffer([]byte{})
		want := fmt.Sprintf("MR: %s\nCommented by: %s\n%s\n\n", mergeRequest.Title, note.Author, mergeRequest.URL)

		app := NewApp(gitlab)
		state, err := app.Run(state, out)
		require.NoError(t, err)
		require.Equal(t, want, out.String())

		out.Reset()
		want = ""

		_, err = app.Run(state, out)
		require.NoError(t, err)
		require.Equal(t, want, out.String())
	})
}
