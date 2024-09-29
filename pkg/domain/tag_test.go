package domain_test

import (
	"msgtm/pkg/domain"
	"reflect"
	"testing"
)

func TestServiceTagUpdate(t *testing.T) {
	tests := []struct {
		name        string
		serviceTag  domain.ServiceTagWithSemVer
		updateMajor bool
		updateMinor bool
		updatePatch bool
		want        domain.ServiceTagWithSemVer
	}{
		{
			name:        "update major",
			serviceTag:  *domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(1, 2, 3)),
			updateMajor: true,
			want:        *domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(2, 0, 0)),
		},
		{
			name:        "update minor",
			serviceTag:  *domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(1, 2, 3)),
			updateMinor: true,
			want:        *domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(1, 3, 0)),
		},
		{
			name:        "update patch",
			serviceTag:  *domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(1, 2, 3)),
			updatePatch: true,
			want:        *domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(1, 2, 4)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.updatePatch {
				tt.serviceTag.UpdatePatch()
			}
			if tt.updateMinor {
				tt.serviceTag.UpdateMinor()
			}
			if tt.updateMajor {
				tt.serviceTag.UpdateMajor()
			}
			if tt.serviceTag.String() != tt.want.String() {
				t.Errorf("ServiceTagUpdate() = %v, want %v", tt.serviceTag, tt.want)
			}
		})
	}

}

func TestSemVerFromStr(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  domain.SemVer
		isErr bool
	}{
		{
			name:  "valid semver string",
			input: "v1.2.3",
			want:  domain.NewSemVer(1, 2, 3),
		},
		{
			name:  "valid semver string without v",
			input: "1.2.3",
			want:  domain.NewSemVer(1, 2, 3),
		},
		{
			name:  "valid semver string v0.0.0",
			input: "v0.0.0",
			want:  domain.NewSemVer(0, 0, 0),
		},
		{
			name:  "valid case add trim string",
			input: "v0.0.0\n",
			want:  domain.NewSemVer(0, 0, 0),
		},
		{
			name:  "invalid semver string",
			input: "v1.2",
			isErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.FromStr(tt.input)
			if (err != nil) != tt.isErr {
				t.Errorf("FromStr() error = %v, wantErr %v", err, tt.isErr)
				return
			}
			if tt.isErr {
				return
			}
			if got.String() != tt.want.String() {
				t.Errorf("FromStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionCompare(t *testing.T) {
	const (
		Greater = iota
		Less
		Equal
	)
	cmpTestCases := []struct {
		name      string
		a         domain.SemVer
		b         domain.SemVer
		aXXXThanB int
	}{
		{
			name:      "major a > b then a greater than b",
			a:         domain.NewSemVer(1, 2, 3),
			b:         domain.NewSemVer(0, 9, 9),
			aXXXThanB: Greater,
		},
		{
			name:      "major a = b and minor a > b then a greater than b",
			a:         domain.NewSemVer(2, 3, 1),
			b:         domain.NewSemVer(2, 2, 3),
			aXXXThanB: Greater,
		},
		{
			name:      "major a = b and minor a = b and patch a > b then a greater than b",
			a:         domain.NewSemVer(1, 2, 4),
			b:         domain.NewSemVer(1, 2, 3),
			aXXXThanB: Greater,
		},
		{
			name:      "major a < b then a less than b",
			a:         domain.NewSemVer(1, 2, 3),
			b:         domain.NewSemVer(2, 0, 0),
			aXXXThanB: Less,
		},
		{
			name:      "major a = b and minor a < b then a less than b",
			a:         domain.NewSemVer(2, 3, 1),
			b:         domain.NewSemVer(2, 4, 3),
			aXXXThanB: Less,
		},
		{
			name:      "major a = b and minor a = b and patch a < b then a less than b",
			a:         domain.NewSemVer(1, 2, 3),
			b:         domain.NewSemVer(1, 2, 4),
			aXXXThanB: Less,
		},
		{
			name:      "a = b",
			a:         domain.NewSemVer(1, 2, 3),
			b:         domain.NewSemVer(1, 2, 3),
			aXXXThanB: Equal,
		},
	}
	for _, tt := range cmpTestCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.aXXXThanB == Greater {
				if !tt.a.GreaterThan(tt.b) {
					t.Errorf("GreaterThan() = %v, want %v", false, true)
				}
			}
			if tt.aXXXThanB == Less {
				if !tt.a.LessThan(tt.b) {
					t.Errorf("LessThan() = %v, want %v", false, true)
				}
			}
			if tt.aXXXThanB == Equal {
				if !tt.a.Equal(tt.b) {
					t.Errorf("Equal() = %v, want %v", false, true)
				}
			}
		})
	}
}

func TestGitTagToServiceTag(t *testing.T) {
	tests := []struct {
		name   string
		gitTag domain.GitTag
		want   domain.ServiceTagWithSemVer
		isErr  bool
	}{

		{
			name:   "valid semver string",
			gitTag: domain.GitTag("service-a-v1.2.3"),
			want:   *domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(1, 2, 3)),
		},
		{
			name:   "valid semver string without v",
			gitTag: domain.GitTag("service-a-1.2.3"),
			want:   *domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(1, 2, 3)),
		},
		{
			name:   "invalid semver string",
			gitTag: domain.GitTag("service-a-v1.2"),
			isErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.gitTag.ToServiceTag()
			if (err != nil) != tt.isErr {
				t.Errorf("FromStrToServiceTag() error = %v, wantErr %v", err, tt.isErr)
				return
			}
			if tt.isErr {
				return
			}
			if got.String() != tt.want.String() {
				t.Errorf("FromStrToServiceTag() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestMajorUpAll(t *testing.T) {
	tests := []struct {
		name    string
		allTags *[]domain.GitTag
		want    *[]*domain.ServiceTagWithSemVer
	}{
		{
			name: "only service tags",
			allTags: &[]domain.GitTag{
				domain.GitTag("service-a-v1.2.3"),
				domain.GitTag("service-b-v2.3.0"),
			},
			want: &[]*domain.ServiceTagWithSemVer{
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(2, 0, 0)),
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-b"), domain.NewSemVer(3, 0, 0)),
			},
		},
		{
			name: "only service tags and duplicate prev service tags",
			allTags: &[]domain.GitTag{
				domain.GitTag("service-a-v1.2.3"),
				domain.GitTag("service-b-v2.3.0"),
				// prev service-b tag
				domain.GitTag("service-b-v2.2.0"),
			},
			want: &[]*domain.ServiceTagWithSemVer{
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(2, 0, 0)),
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-b"), domain.NewSemVer(3, 0, 0)),
			},
		},
		{
			name: "only normal tags",
			allTags: &[]domain.GitTag{
				domain.GitTag("normal-tag"),
			},
			want: &[]*domain.ServiceTagWithSemVer{},
		},
		{
			name: "normal tag and service tags",
			allTags: &[]domain.GitTag{
				domain.GitTag("service-a-v1.2.3"),
				domain.GitTag("normal-tag"),
				domain.GitTag("service-b-v2.3.0"),
			},
			want: &[]*domain.ServiceTagWithSemVer{
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(2, 0, 0)),
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-b"), domain.NewSemVer(3, 0, 0)),
			},
		},
		{
			name:    "no tags",
			allTags: nil,
			want:    &[]*domain.ServiceTagWithSemVer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.MajorUpAll(tt.allTags)
			if !cmpArrayContent(*got, *tt.want) {
				t.Errorf("MajorUpAll() = %v, want %v", got, tt.want)
			}
		})

	}
}

func TestMinorUpAll(t *testing.T) {
	tests := []struct {
		name    string
		allTags *[]domain.GitTag
		want    *[]*domain.ServiceTagWithSemVer
	}{
		{
			name: "only service tags",
			allTags: &[]domain.GitTag{
				domain.GitTag("service-a-v1.2.3"),
				domain.GitTag("service-b-v2.3.0"),
			},

			want: &[]*domain.ServiceTagWithSemVer{
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(1, 3, 0)),
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-b"), domain.NewSemVer(2, 4, 0)),
			},
		},
		{
			name: "only service tags and duplicate prev service tags",
			allTags: &[]domain.GitTag{
				domain.GitTag("service-a-v1.2.3"),
				domain.GitTag("service-b-v2.3.1"),
				// prev service-b tag
				domain.GitTag("service-b-v2.3.0"),
			},
			want: &[]*domain.ServiceTagWithSemVer{
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(1, 3, 0)),
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-b"), domain.NewSemVer(2, 4, 0)),
			},
		},
		{
			name: "normal tag and service tags",
			allTags: &[]domain.GitTag{
				domain.GitTag("service-a-v1.2.3"),
				domain.GitTag("normal-tag"),
				domain.GitTag("service-b-v2.3.0"),
			},

			want: &[]*domain.ServiceTagWithSemVer{
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(1, 3, 0)),
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-b"), domain.NewSemVer(2, 4, 0)),
			},
		},
		{
			name: "only normal tags",
			allTags: &[]domain.GitTag{
				domain.GitTag("normal-tag"),
			},
			want: &[]*domain.ServiceTagWithSemVer{},
		},
		{
			name:    "no tags",
			allTags: nil,
			want:    &[]*domain.ServiceTagWithSemVer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.MinorUpAll(tt.allTags)
			if !cmpArrayContent(*got, *tt.want) {
				t.Errorf("MinorUpAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatchUpAll(t *testing.T) {
	tests := []struct {
		name    string
		allTags *[]domain.GitTag
		want    *[]*domain.ServiceTagWithSemVer
	}{
		{
			name: "only service tags",
			allTags: &[]domain.GitTag{
				domain.GitTag("service-a-v1.2.3"),
				domain.GitTag("service-b-v2.3.0"),
			},

			want: &[]*domain.ServiceTagWithSemVer{
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(1, 2, 4)),
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-b"), domain.NewSemVer(2, 3, 1)),
			},
		},
		{
			name: "only service tags and duplicate prev service tags",
			allTags: &[]domain.GitTag{
				domain.GitTag("service-a-v1.2.3"),
				domain.GitTag("service-b-v2.3.1"),
				// prev service-b tag
				domain.GitTag("service-b-v2.3.0"),
			},
			want: &[]*domain.ServiceTagWithSemVer{
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(1, 2, 4)),
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-b"), domain.NewSemVer(2, 3, 2)),
			},
		},
		{
			name: "normal tag and service tags",
			allTags: &[]domain.GitTag{
				domain.GitTag("service-a-v1.2.3"),
				domain.GitTag("normal-tag"),
				domain.GitTag("service-b-v2.3.0"),
			},

			want: &[]*domain.ServiceTagWithSemVer{
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-a"), domain.NewSemVer(1, 2, 4)),
				domain.NewServiceTagWithSemVer(domain.ServiceName("service-b"), domain.NewSemVer(2, 3, 1)),
			},
		},
		{
			name: "only normal tags",
			allTags: &[]domain.GitTag{
				domain.GitTag("normal-tag"),
			},
			want: &[]*domain.ServiceTagWithSemVer{},
		},
		{
			name:    "no tags",
			allTags: nil,
			want:    &[]*domain.ServiceTagWithSemVer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.PatchUpAll(tt.allTags)
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
