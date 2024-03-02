package main

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type MergeRequest struct {
	ProjectID int
	IID       int
	URL       string
	Title     string
}

type Note struct {
	ID     int
	Author string
}

type XanzyGitlab struct {
	client     *gitlab.Client
	userID     int
	projectIDs []int
}

func NewXanzyGitlab(client *gitlab.Client, userID int, projectIDs []int) *XanzyGitlab {
	return &XanzyGitlab{
		client:     client,
		userID:     userID,
		projectIDs: projectIDs,
	}
}

func (x *XanzyGitlab) ListOpenedMergeRequests() ([]MergeRequest, error) {
	state := "opened"
	scope := "all"

	result := make([]MergeRequest, 0)
	for _, projectID := range x.projectIDs {
		mrs, _, err := x.client.MergeRequests.ListProjectMergeRequests(projectID, &gitlab.ListProjectMergeRequestsOptions{
			State: &state,
			Scope: &scope,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to fetch merge request for %d project: %w", projectID, err)
		}

		for _, mr := range mrs {
			result = append(result, MergeRequest{
				ProjectID: mr.ProjectID,
				IID:       mr.IID,
				URL:       mr.WebURL,
				Title:     mr.Title,
			})
		}
	}

	return result, nil
}

func (x *XanzyGitlab) ListMergeRequestNotes(mr MergeRequest) ([]Note, error) {
	notes, _, err := x.client.Notes.ListMergeRequestNotes(mr.ProjectID, mr.IID, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch notes for %d merge request: %w", mr.IID, err)
	}

	result := make([]Note, 0, len(notes))
	for _, note := range notes {
		if note.System || note.Author.ID == x.userID {
			continue
		}

		result = append(result, Note{
			ID:     note.ID,
			Author: note.Author.Name,
		})
	}

	return result, nil
}
