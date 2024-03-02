package main

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type Gitlab interface {
	ListOpenedMergeRequests() ([]MergeRequest, error)
	ListMergeRequestNotes(mr MergeRequest) ([]Note, error)
}

type NoteID string

func NewNoteID(mr MergeRequest, note Note) NoteID {
	return NoteID(fmt.Sprintf("%d:%d:%d", mr.ProjectID, mr.IID, note.ID))
}

type State struct {
	Notes map[NoteID]struct{} `json:"notes"`
}

func NewState() State {
	return State{
		Notes: make(map[NoteID]struct{}),
	}
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
			noteID := NewNoteID(mergeRequest, note)
			newState.Notes[noteID] = struct{}{}
			if _, ok := state.Notes[noteID]; ok {
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
