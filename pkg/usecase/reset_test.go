package usecase_test

import (
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
	"testing"
)

func TestResetTags(t *testing.T) {
	commitGetter := &StubCommitGetter{
		commitId: domain.HEAD,
		tags: []domain.GitTag{
			domain.GitTag("service-a-v1.2.3"),
			domain.GitTag("service-b-v1.2.3"),
		},
	}

	mockDestroyer := &MockDestroyer{}
	// commitと同じタグを全て削除する
	h := domain.HEAD
	err := usecase.ResetServiceTags(mockDestroyer, commitGetter, &h)
	if err != nil {
		t.Errorf("ResetTags() error = %v, want nil", err)
	}
	expected := []*domain.ServiceTagWithSemVer{
		domain.NewServiceTagWithSemVer("service-a", domain.NewSemVer(1, 2, 3)),
		domain.NewServiceTagWithSemVer("service-b", domain.NewSemVer(1, 2, 3)),
	}
	if !cmpArrayContent(
		*mockDestroyer.Destroyed,
		expected,
	) {
		t.Errorf("ResetTags() = %v, want %v", mockDestroyer.Destroyed, []domain.ServiceTagWithSemVer{})
	}
}
