package main

import (
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
	userId     int
	projectIds []int
}

func NewXanzyGitlab(client *gitlab.Client, userId int, projectIds []int) *XanzyGitlab {
	return &XanzyGitlab{
		client:     client,
		userId:     userId,
		projectIds: projectIds,
	}
}

func (x *XanzyGitlab) ListOpenedMergeRequests() ([]MergeRequest, error) {
	state := "opened"
	scope := "all"

	result := make([]MergeRequest, 0)
	for _, projectId := range x.projectIds {
		mrs, _, err := x.client.MergeRequests.ListProjectMergeRequests(projectId, &gitlab.ListProjectMergeRequestsOptions{
			State: &state,
			Scope: &scope,
		})
		if err != nil {
			return nil, err
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
		return nil, err
	}

	result := make([]Note, 0, len(notes))
	for _, note := range notes {
		if note.System || note.Author.ID == x.userId {
			continue
		}

		result = append(result, Note{
			ID:     note.ID,
			Author: note.Author.Name,
		})
	}
	return result, nil
}
