package main

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type Gitlab interface {
	ListOpenedMergeRequests() ([]MergeRequest, error)
	ListMergeRequestNotes(mergeRequest MergeRequest) ([]Note, error)
}

type State struct {
	Notes map[string]struct{} `json:"notes"`
}

func NewState() State {
	return State{
		Notes: make(map[string]struct{}),
	}
}

func (s *State) HasNote(mergeRequest MergeRequest, note Note) bool {
	key := s.createNoteKey(mergeRequest, note)
	_, ok := s.Notes[key]

	return ok
}

func (s *State) AddNote(mergeRequest MergeRequest, note Note) {
	key := s.createNoteKey(mergeRequest, note)
	s.Notes[key] = struct{}{}
}

func (s *State) createNoteKey(mergeRequest MergeRequest, note Note) string {
	return fmt.Sprintf("%d:%d:%d", mergeRequest.ProjectID, mergeRequest.IID, note.ID)
}

type App struct {
	gitlab Gitlab
}

func NewApp(gitlab Gitlab) *App {
	return &App{gitlab: gitlab}
}

func (a *App) Run(state State, out io.Writer) (State, error) {
	mergeRequests, err := a.gitlab.ListOpenedMergeRequests()
	if err != nil {
		return state, fmt.Errorf("unable to fetch merge requests: %w", err)
	}

	newState := NewState()
	for _, mergeRequest := range mergeRequests {
		notes, err := a.gitlab.ListMergeRequestNotes(mergeRequest)
		if err != nil {
			return state, fmt.Errorf("unable to load notes for merge request: %w", err)
		}

		newCommentsBy := make([]string, 0, len(notes))
		for _, note := range notes {
			newState.AddNote(mergeRequest, note)
			if state.HasNote(mergeRequest, note) {
				continue
			}

			newCommentsBy = append(newCommentsBy, note.Author)
		}

		newCommentsBy = slices.Compact(newCommentsBy)
		slices.Sort(newCommentsBy)

		if len(newCommentsBy) > 0 {
			fmt.Fprintln(out, "MR:", mergeRequest.Title)
			fmt.Fprintln(out, "Commented by:", strings.Join(newCommentsBy, ", "))
			fmt.Fprintln(out, mergeRequest.URL)
			fmt.Fprintln(out, "")
		}
	}

	return newState, nil
}
