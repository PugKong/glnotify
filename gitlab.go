package main

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type MergeRequest struct {
	ProjectID int      `json:"project_id"`
	IID       int      `json:"iid"`
	URL       string   `json:"url"`
	Title     string   `json:"title"`
	Labels    []string `json:"labels"`
}

type Note struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
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
		mergeRequests, _, err := x.client.MergeRequests.ListProjectMergeRequests(
			projectID,
			&gitlab.ListProjectMergeRequestsOptions{
				State: &state,
				Scope: &scope,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch merge request for %d project: %w", projectID, err)
		}

		for _, mergeRequest := range mergeRequests {
			result = append(result, MergeRequest{
				ProjectID: mergeRequest.ProjectID,
				IID:       mergeRequest.IID,
				URL:       mergeRequest.WebURL,
				Title:     mergeRequest.Title,
				Labels:    mergeRequest.Labels,
			})
		}
	}

	return result, nil
}

func (x *XanzyGitlab) ListMergeRequestNotes(mergeRequest MergeRequest) ([]Note, error) {
	notes, _, err := x.client.Notes.ListMergeRequestNotes(mergeRequest.ProjectID, mergeRequest.IID, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch notes for %d merge request: %w", mergeRequest.IID, err)
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
