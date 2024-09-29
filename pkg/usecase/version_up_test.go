package usecase_test

import (
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
	"testing"
)

func TestMajorVersionUpAll(t *testing.T) {
	stub := &StubTagList{
		tags: &[]domain.GitTag{
			domain.GitTag("service-a-v1.2.3"),
			domain.GitTag("service-b-v1.2.3"),
			// prev version
			domain.GitTag("service-a-v1.2.2"),
			// prev version
			domain.GitTag("service-a-v0.2.2"),
		},
	}
	mockRegister := &MockRegister{}
	h := domain.HEAD
	err := usecase.VersionUpAllServiceTags(stub, mockRegister, domain.MajorUpAll, &h)
	if err != nil {
		t.Errorf("MajorUpAllServiceTags() error = %v, want nil", err)
	}
	expected := []*domain.ServiceTagWithSemVer{
		domain.NewServiceTagWithSemVer("service-a", domain.NewSemVer(2, 0, 0)),
		domain.NewServiceTagWithSemVer("service-b", domain.NewSemVer(2, 0, 0)),
	}
	if !cmpArrayContent(
		*mockRegister.AddedTags,
		expected,
	) {
		t.Errorf(
			"MajorUpAllServiceTags() = %v, want %v", mockRegister.AddedTags, expected,
		)
	}
}

func TestMinorVersionUpAll(t *testing.T) {
	stub := &StubTagList{
		tags: &[]domain.GitTag{
			domain.GitTag("service-a-v1.2.3"),
			domain.GitTag("service-b-v1.3.3"),
			// prev version
			domain.GitTag("service-a-v1.2.2"),
		},
	}
	mockRegister := &MockRegister{}
	h := domain.HEAD
	err := usecase.VersionUpAllServiceTags(stub, mockRegister, domain.MinorUpAll, &h)
	if err != nil {
		t.Errorf("MinorUpAllServiceTags() error = %v, want nil", err)
	}
	expected := []*domain.ServiceTagWithSemVer{
		domain.NewServiceTagWithSemVer("service-a", domain.NewSemVer(1, 3, 0)),
		domain.NewServiceTagWithSemVer("service-b", domain.NewSemVer(1, 4, 0)),
	}
	if !cmpArrayContent(
		*mockRegister.AddedTags,
		expected,
	) {
		t.Errorf(
			"MinorUpAllServiceTags() = %v, want %v", mockRegister.AddedTags, expected,
		)
	}
}

func TestPatchVersionUpAll(t *testing.T) {
	stub := &StubTagList{
		tags: &[]domain.GitTag{
			domain.GitTag("service-a-v1.2.3"),
			domain.GitTag("service-b-v1.3.3"),
			// prev version
			domain.GitTag("service-a-v1.2.2"),
		},
	}
	mockRegister := &MockRegister{}
	h := domain.HEAD
	err := usecase.VersionUpAllServiceTags(stub, mockRegister, domain.PatchUpAll, &h)
	if err != nil {
		t.Errorf("PatchUpAllServiceTags() error = %v, want nil", err)
	}
	expected := []*domain.ServiceTagWithSemVer{
		domain.NewServiceTagWithSemVer("service-a", domain.NewSemVer(1, 2, 4)),
		domain.NewServiceTagWithSemVer("service-b", domain.NewSemVer(1, 3, 4)),
	}

	if !cmpArrayContent(
		*mockRegister.AddedTags,
		expected,
	) {
		t.Errorf(
			"PatchUpAllServiceTags() = %v, want %v", mockRegister.AddedTags, expected,
		)
	}
}
