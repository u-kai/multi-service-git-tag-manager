package msgtm_test

import (
	"msgtm"
	"testing"
)

type MockGit struct {
	AddedTags []*msgtm.GitTag
}

func (m *MockGit) AddTag(tag *msgtm.GitTag, commitId *msgtm.CommitId) error {
	m.AddedTags = append(m.AddedTags, tag)
	return nil
}

func TestMajorVersionUpAll(t *testing.T) {
	stub := &StubTagList{
		tags: &[]msgtm.GitTag{
			msgtm.GitTag("service-a-v1.2.3"),
			msgtm.GitTag("service-b-v1.2.3"),
			// prev version
			msgtm.GitTag("service-a-v1.2.2"),
		},
	}
	mockGit := &MockGit{}
	h := msgtm.HEAD
	err := msgtm.VersionUpAllServiceTags(stub, mockGit, msgtm.MajorUpAll, &h)
	if err != nil {
		t.Errorf("MajorUpAllServiceTags() error = %v, want nil", err)
	}
	tag1 := msgtm.GitTag("service-a-v2.0.0")
	tag2 := msgtm.GitTag("service-b-v2.0.0")
	expected := []*msgtm.GitTag{
		&tag1,
		&tag2,
	}
	if !cmpArrayContent(
		mockGit.AddedTags,
		expected,
	) {
		t.Errorf(
			"MajorUpAllServiceTags() = %v, want %v", mockGit.AddedTags, expected,
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
	mockGit := &MockGit{}
	h := msgtm.HEAD
	err := msgtm.VersionUpAllServiceTags(stub, mockGit, msgtm.MinorUpAll, &h)
	if err != nil {
		t.Errorf("MinorUpAllServiceTags() error = %v, want nil", err)
	}
	tag1 := msgtm.GitTag("service-a-v1.3.0")
	tag2 := msgtm.GitTag("service-b-v1.4.0")
	expected := []*msgtm.GitTag{
		&tag1,
		&tag2,
	}
	if !cmpArrayContent(
		mockGit.AddedTags,
		expected,
	) {
		t.Errorf(
			"MinorUpAllServiceTags() = %v, want %v", mockGit.AddedTags, expected,
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
	mockGit := &MockGit{}
	h := msgtm.HEAD
	err := msgtm.VersionUpAllServiceTags(stub, mockGit, msgtm.PatchUpAll, &h)
	if err != nil {
		t.Errorf("PatchUpAllServiceTags() error = %v, want nil", err)
	}
	tag1 := msgtm.GitTag("service-a-v1.2.4")
	tag2 := msgtm.GitTag("service-b-v1.3.4")
	expected := []*msgtm.GitTag{
		&tag1,
		&tag2,
	}
	if !cmpArrayContent(
		mockGit.AddedTags,
		expected,
	) {
		t.Errorf(
			"PatchUpAllServiceTags() = %v, want %v", mockGit.AddedTags, expected,
		)
	}
}
