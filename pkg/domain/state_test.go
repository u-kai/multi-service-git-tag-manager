package domain_test

import (
	"encoding/json"
	"msgtm/pkg/domain"
	"reflect"
	"testing"
)

func newServiceName(name string) *domain.ServiceName {
	s := domain.ServiceName(name)
	return &s
}
func newCommitId(id string) *domain.CommitId {
	c := domain.CommitId(id)
	return &c
}
func newPtrString(s string) *string {
	return &s
}
func TestUpdate(t *testing.T) {
	tests := []struct {
		name        string
		state       *domain.WritedState
		updateInfos []domain.ServiceTagInfo
		want        domain.WritedState
	}{
		{
			name:  "from init",
			state: &domain.WritedState{},
			updateInfos: []domain.ServiceTagInfo{
				{
					Tag:      domain.NewServiceTagWithSemVer(*newServiceName("service1"), domain.SemVer{Major: 1, Minor: 0, Patch: 0}),
					CommitId: newCommitId("commit1"),
				},
				{
					Tag:      domain.NewServiceTagWithSemVer(*newServiceName("service2"), domain.SemVer{Major: 1, Minor: 0, Patch: 0}),
					CommitId: newCommitId("commit2"),
				},
			},
			want: domain.WritedState{
				ServiceTagStates: []*domain.ServiceTagState{
					{
						ServiceName: newServiceName("service1"),
						Latest: &domain.ServiceTagInfo{
							Tag:      domain.NewServiceTagWithSemVer(*newServiceName("service1"), domain.SemVer{Major: 1, Minor: 0, Patch: 0}),
							CommitId: newCommitId("commit1"),
						},
					},
					{
						ServiceName: newServiceName("service2"),
						Latest: &domain.ServiceTagInfo{
							Tag:      domain.NewServiceTagWithSemVer(*newServiceName("service2"), domain.SemVer{Major: 1, Minor: 0, Patch: 0}),
							CommitId: newCommitId("commit2"),
						},
					},
				},
			},
		},
		{
			name:  "update state",
			state: domain.InitStateWriter("service1", "service2"),
			updateInfos: []domain.ServiceTagInfo{
				{
					Tag:      domain.NewServiceTagWithSemVer(*newServiceName("service1"), domain.SemVer{Major: 1, Minor: 0, Patch: 0}),
					CommitId: newCommitId("commit1"),
				},
				{
					Tag:      domain.NewServiceTagWithSemVer(*newServiceName("service2"), domain.SemVer{Major: 1, Minor: 0, Patch: 0}),
					CommitId: newCommitId("commit2"),
				},
			},
			want: domain.WritedState{
				ServiceTagStates: []*domain.ServiceTagState{
					{
						ServiceName: newServiceName("service1"),
						Latest: &domain.ServiceTagInfo{
							Tag:      domain.NewServiceTagWithSemVer(*newServiceName("service1"), domain.SemVer{Major: 1, Minor: 0, Patch: 0}),
							CommitId: newCommitId("commit1"),
						},
						Prev: nil,
					},
					{
						ServiceName: newServiceName("service2"),
						Latest: &domain.ServiceTagInfo{
							Tag:      domain.NewServiceTagWithSemVer(*newServiceName("service2"), domain.SemVer{Major: 1, Minor: 0, Patch: 0}),
							CommitId: newCommitId("commit2"),
						},
					},
				},
			},
		},
		{
			name:        "nil update state",
			state:       domain.InitStateWriter("service1", "service2"),
			updateInfos: nil,
			want: domain.WritedState{
				ServiceTagStates: []*domain.ServiceTagState{
					{
						ServiceName: newServiceName("service1"),
						Latest:      nil,
						Prev:        nil,
					},
					{
						ServiceName: newServiceName("service2"),
						Latest:      nil,
						Prev:        nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, info := range tt.updateInfos {
				tt.state.Update(info.Tag.Service, &info)
			}
			if !reflect.DeepEqual(tt.state, &tt.want) {
				t.Errorf("got: %v, want: %v", tt.state, tt.want)
				jsonGot, _ := json.Marshal(tt.state)
				jsonWant, _ := json.Marshal(tt.want)
				t.Errorf("got: %s, want: %s", jsonGot, jsonWant)
			}
		})
	}
}
