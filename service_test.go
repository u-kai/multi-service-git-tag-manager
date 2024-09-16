package msgtm_test

import (
	"msgtm"
	"reflect"
	"testing"
)

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
				msgtm.GitTag("service-b-v2.2.0"),
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
			got := msgtm.MajorUpAll(tt.allTags)
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
			name: "only service tags and duplicate prev service tags",
			allTags: &[]msgtm.GitTag{
				msgtm.GitTag("service-a-v1.2.3"),
				msgtm.GitTag("service-b-v2.3.1"),
				// prev service-b tag
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
			got := msgtm.MinorUpAll(tt.allTags)
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
			name: "only service tags and duplicate prev service tags",
			allTags: &[]msgtm.GitTag{
				msgtm.GitTag("service-a-v1.2.3"),
				msgtm.GitTag("service-b-v2.3.1"),
				// prev service-b tag
				msgtm.GitTag("service-b-v2.3.0"),
			},
			want: &[]*msgtm.ServiceTagWithSemVer{
				msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 2, 4)),
				msgtm.NewServiceTagWithSemVer("service-b", msgtm.NewSemVer(2, 3, 2)),
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
			got := msgtm.PatchUpAll(tt.allTags)
			if !cmpArrayContent(*got, *tt.want) {
				t.Errorf("PatchUpAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 順不同な配列の比較
func cmpArrayContent[T any](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		found := false
		for _, vv := range b {
			if reflect.DeepEqual(v, vv) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
