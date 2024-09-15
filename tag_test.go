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
func TestFromStrToServiceTag(t *testing.T) {
	tests := []struct {
		name  string
		str   string
		want  msgtm.ServiceTagWithSemVer
		isErr bool
	}{

		{
			name: "valid semver string",
			str:  "service-a-v1.2.3",
			want: *msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 2, 3)),
		},
		{
			name: "valid semver string without v",
			str:  "service-a-1.2.3",
			want: *msgtm.NewServiceTagWithSemVer("service-a", msgtm.NewSemVer(1, 2, 3)),
		},
		{
			name:  "invalid semver string",
			str:   "service-a-v1.2",
			isErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := msgtm.FromStrToServiceTag(tt.str)
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
