package main

import (
	"fmt"
	"slices"
	"strings"
)

type Gitlab interface {
	ListOpenedMergeRequests() ([]MergeRequest, error)
	ListMergeRequestNotes(mr MergeRequest) ([]Note, error)
}

type NoteID string

func NewNoteId(mr MergeRequest, note Note) NoteID {
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

func (a *App) Run(state State) (State, error) {
	mrs, err := a.gitlab.ListOpenedMergeRequests()
	if err != nil {
		return state, err
	}

	newState := NewState()
	for _, mr := range mrs {
		notes, err := a.gitlab.ListMergeRequestNotes(mr)
		if err != nil {
			return state, err
		}

		newCommentsBy := make([]string, 0, len(notes))
		for _, note := range notes {
			noteId := NewNoteId(mr, note)
			newState.Notes[noteId] = struct{}{}
			if _, ok := state.Notes[noteId]; ok {
				continue
			}

			newCommentsBy = append(newCommentsBy, note.Author)
		}

		newCommentsBy = slices.Compact(newCommentsBy)
		slices.Sort(newCommentsBy)

		if len(newCommentsBy) > 0 {
			fmt.Println("MR:", mr.Title)
			fmt.Println("Commented by:", strings.Join(newCommentsBy, ", "))
			fmt.Println(mr.URL)
			fmt.Println("")
		}
	}

	return newState, nil
}
