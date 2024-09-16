package msgtm_test

import (
	"msgtm"
	"reflect"
	"testing"
)

type StubTagList struct {
	tags *[]msgtm.GitTag
}

func (s *StubTagList) List() (*[]msgtm.GitTag, error) {
	return s.tags, nil
}

func TestMinorUpAll(t *testing.T) {
	tests := []struct {
		name    string
		allTags *[]msgtm.GitTag
		want    *[]*msgtm.ServiceTagWithSemVer
	}{
		{
			name: "only service tags",
			allTags: &[]msgtm.GitTag{
				msgtm.GitTag("service-a-v1.2.3"),
				msgtm.GitTag("service-b-v2.3.0"),
			},

			want: &[]*msgtm.ServiceTagWithSemVer{
				msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 3, 0)),
				msgtm.NewServiceTagWithSemVer("service-b", msgtm.NewSemVer(2, 4, 0)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stubTagList := &StubTagList{tags: tt.allTags}
			got, err := msgtm.MinorUpAll(stubTagList)
			if err != nil {
				t.Errorf("MinorUpAll() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MinorUpAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
