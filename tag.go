package msgtm

import (
	"fmt"
	"regexp"
)

type SemVer struct {
	major int
	minor int
	patch int
}

func FromStr(s string) (SemVer, error) {
	var major, minor, patch int
	_, err := fmt.Sscanf(s, "v%d.%d.%d", &major, &minor, &patch)
	if err != nil {
		_, err = fmt.Sscanf(s, "%d.%d.%d", &major, &minor, &patch)
		if err != nil {
			return SemVer{}, fmt.Errorf("invalid semver string: %s", s)
		}
	}
	return NewSemVer(major, minor, patch), nil

}
func NewSemVer(major, minor, patch int) SemVer {
	return SemVer{
		major: major,
		minor: minor,
		patch: patch,
	}
}

func (s SemVer) MajorUp() SemVer {
	return SemVer{
		major: s.major + 1,
		minor: 0,
		patch: 0,
	}
}

func (s SemVer) MinorUp() SemVer {
	return SemVer{
		major: s.major,
		minor: s.minor + 1,
		patch: 0,
	}
}

func (s SemVer) PatchUp() SemVer {
	return SemVer{
		major: s.major,
		minor: s.minor,
		patch: s.patch + 1,
	}
}

func (s *SemVer) String() string {
	return fmt.Sprintf("v%d.%d.%d", s.major, s.minor, s.patch)
}

type ServiceTagWithSemVer struct {
	service string
	tag     SemVer
}

var (
	serviceSemVerRe         = regexp.MustCompile(`^([a-zA-Z0-9-]+)-v(\d+)\.(\d+)\.(\d+)$`)
	serviceSemVerWithoutVRe = regexp.MustCompile(`^([a-zA-Z0-9-]+)-(\d+)\.(\d+)\.(\d+)$`)
)

func errInvalidServiceSemVerMsg(invalid string) error {
	return fmt.Errorf("invalid service semver string: %s\nservice tag should be SERVICE_NAME-vMAJOR.MINOR.PATCH", invalid)
}

func FromStrToServiceTag(s string) (*ServiceTagWithSemVer, error) {
	matches := serviceSemVerRe.FindStringSubmatch(s)
	semVerStr := "v"
	if len(matches) != 5 {
		matches = serviceSemVerWithoutVRe.FindStringSubmatch(s)
		if len(matches) != 5 {
			return nil, errInvalidServiceSemVerMsg(s)
		}
		// without v version
		semVerStr = ""
	}

	service := matches[1]
	tag, err := FromStr(fmt.Sprintf("%s%s.%s.%s", semVerStr, matches[2], matches[3], matches[4]))
	if err != nil {
		return nil, errInvalidServiceSemVerMsg(s)
	}

	return NewServiceTagWithSemVer(service, tag), nil
}

func NewServiceTagWithSemVer(service string, tag SemVer) *ServiceTagWithSemVer {
	return &ServiceTagWithSemVer{
		service: service,
		tag:     tag,
	}
}

func (s *ServiceTagWithSemVer) UpdateMajor() {
	s.tag = s.tag.MajorUp()
}
func (s *ServiceTagWithSemVer) UpdateMinor() {
	s.tag = s.tag.MinorUp()
}
func (s *ServiceTagWithSemVer) UpdatePatch() {
	s.tag = s.tag.PatchUp()
}
func (s *ServiceTagWithSemVer) String() string {
	return fmt.Sprintf("%s-%s", s.service, s.tag.String())
}
