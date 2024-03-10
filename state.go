package main

import "fmt"

type MergeRequestState struct {
	MergeRequest MergeRequest `json:"merge_request"`
	Notes        map[int]Note `json:"notes"`
}

func NewMergeRequestState(mergeRequest MergeRequest, notes []Note) MergeRequestState {
	notesMap := make(map[int]Note)
	for _, note := range notes {
		notesMap[note.ID] = note
	}

	return MergeRequestState{
		MergeRequest: mergeRequest,
		Notes:        notesMap,
	}
}

func (m MergeRequestState) HasNote(note Note) bool {
	_, ok := m.Notes[note.ID]

	return ok
}

type State struct {
	MergeRequests map[string]MergeRequestState `json:"merge_requests"`
}

func NewState(mergeRequestStates ...MergeRequestState) State {
	mergeRequestMap := make(map[string]MergeRequestState)
	for _, mergeRequestState := range mergeRequestStates {
		key := createMergeRequestKey(mergeRequestState.MergeRequest)
		mergeRequestMap[key] = mergeRequestState
	}

	return State{
		MergeRequests: mergeRequestMap,
	}
}

func (s State) GetMergeRequestState(mergeRequest MergeRequest) MergeRequestState {
	key := createMergeRequestKey(mergeRequest)

	return s.MergeRequests[key]
}

func createMergeRequestKey(mergeRequest MergeRequest) string {
	return fmt.Sprintf("%d:%d", mergeRequest.ProjectID, mergeRequest.IID)
}
