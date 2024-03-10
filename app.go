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

type App struct {
	out    io.Writer
	gitlab Gitlab
}

func NewApp(out io.Writer, gitlab Gitlab) *App {
	return &App{
		out:    out,
		gitlab: gitlab,
	}
}

func (a *App) Run(prevState State) (State, error) {
	curState, err := a.createNewState()
	if err != nil {
		return curState, err
	}

	for _, curMergeRequestState := range curState.MergeRequests {
		prevMergeRequestState := prevState.GetMergeRequestState(curMergeRequestState.MergeRequest)
		a.printMergeRequest(prevMergeRequestState, curMergeRequestState)
	}

	return curState, nil
}

func (a *App) createNewState() (State, error) {
	mergeRequests, err := a.gitlab.ListOpenedMergeRequests()
	if err != nil {
		return NewState(), fmt.Errorf("unable to fetch merge requests: %w", err)
	}

	mergeRequestStates := make([]MergeRequestState, 0, len(mergeRequests))
	for _, mergeRequest := range mergeRequests {
		notes, err := a.gitlab.ListMergeRequestNotes(mergeRequest)
		if err != nil {
			return NewState(), fmt.Errorf("unable to load notes for merge request: %w", err)
		}

		mergeRequestStates = append(mergeRequestStates, NewMergeRequestState(mergeRequest, notes))
	}

	return NewState(mergeRequestStates...), nil
}

func (a *App) printMergeRequest(prevState, curState MergeRequestState) {
	mergeRequest := curState.MergeRequest
	newLabels := make([]string, 0, len(mergeRequest.Labels))
	for _, label := range mergeRequest.Labels {
		if slices.Contains(prevState.MergeRequest.Labels, label) {
			continue
		}

		newLabels = append(newLabels, label)
	}

	newCommentsBy := make([]string, 0, len(curState.Notes))
	for _, note := range curState.Notes {
		if prevState.HasNote(note) {
			continue
		}

		newCommentsBy = append(newCommentsBy, note.Author)
	}
	newCommentsBy = slices.Compact(newCommentsBy)
	slices.Sort(newCommentsBy)

	if len(newLabels) == 0 && len(newCommentsBy) == 0 {
		return
	}

	fmt.Fprintln(a.out, "MR:", mergeRequest.Title)

	if len(newLabels) > 0 {
		fmt.Fprintln(a.out, "Labeled:", strings.Join(newLabels, ", "))
	}

	if len(newCommentsBy) > 0 {
		fmt.Fprintln(a.out, "Commented by:", strings.Join(newCommentsBy, ", "))
	}

	fmt.Fprintln(a.out, mergeRequest.URL)
	fmt.Fprintln(a.out, "")
}
