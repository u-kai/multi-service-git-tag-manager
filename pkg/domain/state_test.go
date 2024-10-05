package domain_test

import (
	"encoding/json"
	"msgtm/pkg/domain"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
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

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want domain.WritedState
	}{
		{
			name: "valid init service",
			data: []byte(
				`
services:      
    - name: test
      latest: null
      prev: null`),
			want: domain.WritedState{
				ServiceTagStates: []*domain.ServiceTagState{
					{
						ServiceName: newServiceName("test"),
						Latest:      nil,
						Prev:        nil,
					},
				},
			},
		},
		{
			name: "valid service has latest",
			data: []byte(
				`
services:
    - name: test
      latest:
         tag: 
            version: v1.0.0
         commitId: commit1
      prev: null`),
			want: domain.WritedState{
				ServiceTagStates: []*domain.ServiceTagState{
					{
						ServiceName: newServiceName("test"),
						Latest: &domain.ServiceTagInfo{
							Tag:      domain.NewServiceTagWithSemVer(*newServiceName("test"), domain.SemVer{Major: 1, Minor: 0, Patch: 0}),
							CommitId: newCommitId("commit1"),
						},
						Prev: nil,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got domain.WritedState
			err := yaml.Unmarshal(tt.data, &got)
			if err != nil {
				t.Errorf("failed to unmarshal: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				jsonGot, _ := json.Marshal(got)
				jsonWant, _ := json.Marshal(tt.want)
				t.Errorf("got: %s, want: %s", jsonGot, jsonWant)
			}
		})
	}
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
