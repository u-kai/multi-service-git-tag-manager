package msgtm_test

import (
	"msgtm"
	"testing"
)

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
