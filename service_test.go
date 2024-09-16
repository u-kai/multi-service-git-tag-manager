package msgtm_test

import (
	"msgtm"
	"testing"
)

type StubTagList struct {
	tags *[]msgtm.GitTag
}

func (s *StubTagList) List() (*[]msgtm.GitTag, error) {
	return s.tags, nil
}

func TestMajorUpAll(t *testing.T) {
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
				msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(2, 0, 0)),
				msgtm.NewServiceTagWithSemVer("service-b", msgtm.NewSemVer(3, 0, 0)),
			},
		},
		{
			name: "only service tags and duplicate prev service tags",
			allTags: &[]msgtm.GitTag{
				msgtm.GitTag("service-a-v1.2.3"),
				msgtm.GitTag("service-b-v2.3.0"),
				// prev service-b tag
				msgtm.GitTag("service-b-v2.3.0"),
			},
			want: &[]*msgtm.ServiceTagWithSemVer{
				msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(2, 0, 0)),
				msgtm.NewServiceTagWithSemVer("service-b", msgtm.NewSemVer(3, 0, 0)),
			},
		},
		{
			name: "only normal tags",
			allTags: &[]msgtm.GitTag{
				msgtm.GitTag("normal-tag"),
			},
			want: &[]*msgtm.ServiceTagWithSemVer{},
		},
		{
			name: "normal tag and service tags",
			allTags: &[]msgtm.GitTag{
				msgtm.GitTag("service-a-v1.2.3"),
				msgtm.GitTag("normal-tag"),
				msgtm.GitTag("service-b-v2.3.0"),
			},
			want: &[]*msgtm.ServiceTagWithSemVer{
				msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(2, 0, 0)),
				msgtm.NewServiceTagWithSemVer("service-b", msgtm.NewSemVer(3, 0, 0)),
			},
		},
		{
			name:    "no tags",
			allTags: nil,
			want:    &[]*msgtm.ServiceTagWithSemVer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stubTagList := &StubTagList{tags: tt.allTags}
			got, err := msgtm.MajorUpAll(stubTagList)
			if err != nil {
				t.Errorf("MajorUpAll() error = %v", err)
			}
			if !cmpArrayContent(*got, *tt.want) {
				t.Errorf("MajorUpAll() = %v, want %v", got, tt.want)
			}
		})

	}
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
		{
			name: "normal tag and service tags",
			allTags: &[]msgtm.GitTag{
				msgtm.GitTag("service-a-v1.2.3"),
				msgtm.GitTag("normal-tag"),
				msgtm.GitTag("service-b-v2.3.0"),
			},

			want: &[]*msgtm.ServiceTagWithSemVer{
				msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 3, 0)),
				msgtm.NewServiceTagWithSemVer("service-b", msgtm.NewSemVer(2, 4, 0)),
			},
		},
		{
			name: "only normal tags",
			allTags: &[]msgtm.GitTag{
				msgtm.GitTag("normal-tag"),
			},
			want: &[]*msgtm.ServiceTagWithSemVer{},
		},
		{
			name:    "no tags",
			allTags: nil,
			want:    &[]*msgtm.ServiceTagWithSemVer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stubTagList := &StubTagList{tags: tt.allTags}
			got, err := msgtm.MinorUpAll(stubTagList)
			if err != nil {
				t.Errorf("MinorUpAll() error = %v", err)
			}
			if !cmpArrayContent(*got, *tt.want) {
				t.Errorf("MinorUpAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatchUpAll(t *testing.T) {
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
				msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 2, 4)),
				msgtm.NewServiceTagWithSemVer("service-b", msgtm.NewSemVer(2, 3, 1)),
			},
		},
		{
			name: "normal tag and service tags",
			allTags: &[]msgtm.GitTag{
				msgtm.GitTag("service-a-v1.2.3"),
				msgtm.GitTag("normal-tag"),
				msgtm.GitTag("service-b-v2.3.0"),
			},

			want: &[]*msgtm.ServiceTagWithSemVer{
				msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 2, 4)),
				msgtm.NewServiceTagWithSemVer("service-b", msgtm.NewSemVer(2, 3, 1)),
			},
		},
		{
			name: "only normal tags",
			allTags: &[]msgtm.GitTag{
				msgtm.GitTag("normal-tag"),
			},
			want: &[]*msgtm.ServiceTagWithSemVer{},
		},
		{
			name:    "no tags",
			allTags: nil,
			want:    &[]*msgtm.ServiceTagWithSemVer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stubTagList := &StubTagList{tags: tt.allTags}
			got, err := msgtm.PatchUpAll(stubTagList)
			if err != nil {
				t.Errorf("PatchUpAll() error = %v", err)
			}
			if !cmpArrayContent(*got, *tt.want) {
				t.Errorf("PatchUpAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 順不同な配列の比較
func cmpArrayContent(a, b []*msgtm.ServiceTagWithSemVer) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		found := false
		for _, vv := range b {
			if v.String() == vv.String() {
				found = true
				continue
			}
		}
		if !found {
			return false
		}
	}
	return true
}
