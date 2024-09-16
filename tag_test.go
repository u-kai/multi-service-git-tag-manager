package msgtm_test

import (
	"msgtm"
	"testing"
)

func TestServiceTagUpdate(t *testing.T) {
	tests := []struct {
		name        string
		serviceTag  msgtm.ServiceTagWithSemVer
		updateMajor bool
		updateMinor bool
		updatePatch bool
		want        msgtm.ServiceTagWithSemVer
	}{
		{
			name:        "update major",
			serviceTag:  *msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 2, 3)),
			updateMajor: true,
			want:        *msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(2, 0, 0)),
		},
		{
			name:        "update minor",
			serviceTag:  *msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 2, 3)),
			updateMinor: true,
			want:        *msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 3, 0)),
		},
		{
			name:        "update patch",
			serviceTag:  *msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 2, 3)),
			updatePatch: true,
			want:        *msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 2, 4)),
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

func TestVersionCompare(t *testing.T) {
	const (
		Greater = iota
		Less
		Equal
	)
	cmpTestCases := []struct {
		name      string
		a         msgtm.SemVer
		b         msgtm.SemVer
		aXXXThanB int
	}{
		{
			name:      "major a > b then a greater than b",
			a:         msgtm.NewSemVer(1, 2, 3),
			b:         msgtm.NewSemVer(0, 9, 9),
			aXXXThanB: Greater,
		},
		{
			name:      "major a = b and minor a > b then a greater than b",
			a:         msgtm.NewSemVer(2, 3, 1),
			b:         msgtm.NewSemVer(2, 2, 3),
			aXXXThanB: Greater,
		},
		{
			name:      "major a = b and minor a = b and patch a > b then a greater than b",
			a:         msgtm.NewSemVer(1, 2, 4),
			b:         msgtm.NewSemVer(1, 2, 3),
			aXXXThanB: Greater,
		},
		{
			name:      "major a < b then a less than b",
			a:         msgtm.NewSemVer(1, 2, 3),
			b:         msgtm.NewSemVer(2, 0, 0),
			aXXXThanB: Less,
		},
		{
			name:      "major a = b and minor a < b then a less than b",
			a:         msgtm.NewSemVer(2, 3, 1),
			b:         msgtm.NewSemVer(2, 4, 3),
			aXXXThanB: Less,
		},
		{
			name:      "major a = b and minor a = b and patch a < b then a less than b",
			a:         msgtm.NewSemVer(1, 2, 3),
			b:         msgtm.NewSemVer(1, 2, 4),
			aXXXThanB: Less,
		},
		{
			name:      "a = b",
			a:         msgtm.NewSemVer(1, 2, 3),
			b:         msgtm.NewSemVer(1, 2, 3),
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
		gitTag msgtm.GitTag
		want   msgtm.ServiceTagWithSemVer
		isErr  bool
	}{

		{
			name:   "valid semver string",
			gitTag: msgtm.GitTag("service-a-v1.2.3"),
			want:   *msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 2, 3)),
		},
		{
			name:   "valid semver string without v",
			gitTag: msgtm.GitTag("service-a-1.2.3"),
			want:   *msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 2, 3)),
		},
		{
			name:   "invalid semver string",
			gitTag: msgtm.GitTag("service-a-v1.2"),
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
