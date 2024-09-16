package msgtm_test

import (
	"msgtm"
)

type MockGit struct {
	AddedTags []*msgtm.GitTag
}

//func (m *MockGit) AddTag(tag *msgtm.GitTag, commitId *msgtm.CommitId) error {
//	m.AddedTags = append(m.AddedTags, tag)
//	return nil
//}

//func TestMajorUpAllUsecase(t *testing.T) {
//	stub := &StubTagList{
//		tags: &[]msgtm.GitTag{
//			msgtm.GitTag("service-a-v1.2.3"),
//			msgtm.GitTag("service-b-v1.2.3"),
//			// prev version
//			msgtm.GitTag("service-a-v1.2.2"),
//		},
//	}
//	mockGit := &MockGit{}
//
//	err := msgtm.MajorUpAllServiceTags(stub, mockGit)
//	if err != nil {
//		t.Errorf("MajorUpAllServiceTags() error = %v, want nil", err)
//	}
//	if !reflect.DeepEqual(mockGit.AddedTags, []*msgtm.GitTag{
//		&msgtm.GitTag("service-a-v2.0.0"),
//		&msgtm.GitTag("service-b-v2.0.0"),
//	}) {
//		t.Errorf("MajorUpAllServiceTags() = %v, want %v", mockGit.AddedTags, []*msgtm.GitTag{
//			&msgtm.GitTag("service-a-v2.0.0"),
//			&msgtm.GitTag("service-b-v2.0.0"),
//		})
//	}
//}
