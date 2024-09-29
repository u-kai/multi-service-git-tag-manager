package usecase_test

import (
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
)

type StubCommitGetter struct {
	commitId domain.CommitId
	tags     []domain.GitTag
}

func (s *StubCommitGetter) Execute(cmd usecase.GetCommitTagQuery) (*[]domain.GitTag, error) {
	if *cmd.CommitId != s.commitId {
		return nil, nil
	}
	return &s.tags, nil
}

type MockRegister struct {
	AddedTags *[]*domain.ServiceTagWithSemVer
}

func (m *MockRegister) Execute(cmd usecase.RegisterServiceTagsCommand) error {
	m.AddedTags = cmd.Tags
	return nil
}

type MockDestroyer struct {
	Destroyed *[]*domain.ServiceTagWithSemVer
}

func (m *MockDestroyer) Execute(cmd usecase.DestroyServiceTagsCommand) error {
	m.Destroyed = cmd.Tags
	return nil
}

type StubTagList struct {
	tags *[]domain.GitTag
}

func (s *StubTagList) Execute(cmd usecase.ListTagsQuery) (*[]domain.GitTag, error) {
	return s.tags, nil
}
