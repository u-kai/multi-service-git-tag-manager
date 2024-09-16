package msgtm_test

import (
	"msgtm"
	"testing"
)

type MockRegister struct {
	AddedTags *[]*msgtm.ServiceTagWithSemVer
}

func (m *MockRegister) Register(_ *msgtm.CommitId, tags *[]*msgtm.ServiceTagWithSemVer) error {
	m.AddedTags = tags
	return nil
}

type StubTagList struct {
	tags *[]msgtm.GitTag
}

func (s *StubTagList) List() (*[]msgtm.GitTag, error) {
	return s.tags, nil
}

type MockDestroyer struct {
	Destroyed []*msgtm.ServiceTagWithSemVer
}

func (m *MockDestroyer) Destroy(tag []*msgtm.ServiceTagWithSemVer) error {
	m.Destroyed = tag
	return nil
}

type StubCommitGetter struct {
	commitId msgtm.CommitId
	tags     []msgtm.GitTag
}

func (s *StubCommitGetter) GetTags(commitId *msgtm.CommitId) ([]msgtm.GitTag, error) {
	if *commitId != s.commitId {
		return []msgtm.GitTag{}, nil
	}
	return s.tags, nil
}

func TestResetTags(t *testing.T) {
	commitGetter := &StubCommitGetter{
		commitId: msgtm.HEAD,
		tags: []msgtm.GitTag{
			msgtm.GitTag("service-a-v1.2.3"),
			msgtm.GitTag("service-b-v1.2.3"),
		},
	}

	mockDestroyer := &MockDestroyer{}
	// commitと同じタグを全て削除する
	h := msgtm.HEAD
	err := msgtm.ResetServiceTags(mockDestroyer, commitGetter, &h)
	if err != nil {
		t.Errorf("ResetTags() error = %v, want nil", err)
	}
	expected := []*msgtm.ServiceTagWithSemVer{
		msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 2, 3)),
		msgtm.NewServiceTagWithSemVer("service-b", msgtm.NewSemVer(1, 2, 3)),
	}
	if !cmpArrayContent(
		mockDestroyer.Destroyed,
		expected,
	) {
		t.Errorf("ResetTags() = %v, want %v", mockDestroyer.Destroyed, []msgtm.ServiceTagWithSemVer{})
	}
}

func TestCreateServiceTags(t *testing.T) {
	services := []string{"service-a", "service-b"}
	version := msgtm.NewSemVer(0, 0, 1)
	commitId := msgtm.HEAD
	mockRegister := &MockRegister{}
	err := msgtm.CreateServiceTags(mockRegister, &commitId, services, version)
	if err != nil {
		t.Errorf("CreateServiceTags() error = %v, want nil", err)
	}
	expected := []*msgtm.ServiceTagWithSemVer{
		msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(0, 0, 1)),
		msgtm.NewServiceTagWithSemVer("service-b", msgtm.NewSemVer(0, 0, 1)),
	}
	if !cmpArrayContent(
		*mockRegister.AddedTags,
		expected,
	) {
		t.Errorf(
			"CreateServiceTags() = %v, want %v", mockRegister.AddedTags, expected,
		)
	}
}

func TestMajorVersionUpAll(t *testing.T) {
	stub := &StubTagList{
		tags: &[]msgtm.GitTag{
			msgtm.GitTag("service-a-v1.2.3"),
			msgtm.GitTag("service-b-v1.2.3"),
			// prev version
			msgtm.GitTag("service-a-v1.2.2"),
			// prev version
			msgtm.GitTag("service-a-v0.2.2"),
		},
	}
	mockRegister := &MockRegister{}
	h := msgtm.HEAD
	err := msgtm.VersionUpAllServiceTags(stub, mockRegister, msgtm.MajorUpAll, &h)
	if err != nil {
		t.Errorf("MajorUpAllServiceTags() error = %v, want nil", err)
	}
	expected := []*msgtm.ServiceTagWithSemVer{
		msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(2, 0, 0)),
		msgtm.NewServiceTagWithSemVer("service-b", msgtm.NewSemVer(2, 0, 0)),
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
		tags: &[]msgtm.GitTag{
			msgtm.GitTag("service-a-v1.2.3"),
			msgtm.GitTag("service-b-v1.3.3"),
			// prev version
			msgtm.GitTag("service-a-v1.2.2"),
		},
	}
	mockRegister := &MockRegister{}
	h := msgtm.HEAD
	err := msgtm.VersionUpAllServiceTags(stub, mockRegister, msgtm.MinorUpAll, &h)
	if err != nil {
		t.Errorf("MinorUpAllServiceTags() error = %v, want nil", err)
	}
	expected := []*msgtm.ServiceTagWithSemVer{
		msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 3, 0)),
		msgtm.NewServiceTagWithSemVer("service-b", msgtm.NewSemVer(1, 4, 0)),
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
		tags: &[]msgtm.GitTag{
			msgtm.GitTag("service-a-v1.2.3"),
			msgtm.GitTag("service-b-v1.3.3"),
			// prev version
			msgtm.GitTag("service-a-v1.2.2"),
		},
	}
	mockRegister := &MockRegister{}
	h := msgtm.HEAD
	err := msgtm.VersionUpAllServiceTags(stub, mockRegister, msgtm.PatchUpAll, &h)
	if err != nil {
		t.Errorf("PatchUpAllServiceTags() error = %v, want nil", err)
	}
	expected := []*msgtm.ServiceTagWithSemVer{
		msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 2, 4)),
		msgtm.NewServiceTagWithSemVer("service-b", msgtm.NewSemVer(1, 3, 4)),
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
