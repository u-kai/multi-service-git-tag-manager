package msgtm

import (
	"fmt"
	"regexp"
)

type GitTag string

func (g GitTag) String() string {
	return string(g)
}

var (
	serviceSemVerRe         = regexp.MustCompile(`^([a-zA-Z0-9-]+)-v(\d+)\.(\d+)\.(\d+)$`)
	serviceSemVerWithoutVRe = regexp.MustCompile(`^([a-zA-Z0-9-]+)-(\d+)\.(\d+)\.(\d+)$`)
)

func errInvalidServiceSemVerMsg(invalid string) error {
	return fmt.Errorf("invalid service semver string: %s\nservice version should be SERVICE_NAME-vMAJOR.MINOR.PATCH", invalid)
}

func (g GitTag) ToServiceTag() (*ServiceTagWithSemVer, error) {
	matches := serviceSemVerRe.FindStringSubmatch(g.String())
	semVerStr := "v"
	if len(matches) != 5 {
		matches = serviceSemVerWithoutVRe.FindStringSubmatch(g.String())
		if len(matches) != 5 {
			return nil, errInvalidServiceSemVerMsg(g.String())
		}
		// without v version
		semVerStr = ""
	}

	service := matches[1]
	version, err := FromStr(fmt.Sprintf("%s%s.%s.%s", semVerStr, matches[2], matches[3], matches[4]))
	if err != nil {
		return nil, errInvalidServiceSemVerMsg(g.String())
	}

	return NewServiceTagWithSemVer(service, version), nil
}

func FilterServiceTags(tags *[]GitTag) *[]*ServiceTagWithSemVer {
	serviceTags := []*ServiceTagWithSemVer{}
	if tags == nil || len(*tags) == 0 {
		return &serviceTags
	}

	for _, tag := range *tags {
		serviceTag, err := tag.ToServiceTag()
		if err != nil {
			continue
		}
		serviceTags = append(serviceTags, serviceTag)
	}

	return &serviceTags
}

type ServiceTagWithSemVer struct {
	service string
	version SemVer
}

func NewServiceTagWithSemVer(service string, version SemVer) *ServiceTagWithSemVer {
	return &ServiceTagWithSemVer{
		service: service,
		version: version,
	}
}

func (s *ServiceTagWithSemVer) UpdateMajor() {
	s.version = s.version.MajorUp()
}
func (s *ServiceTagWithSemVer) UpdateMinor() {
	s.version = s.version.MinorUp()
}
func (s *ServiceTagWithSemVer) UpdatePatch() {
	s.version = s.version.PatchUp()
}
func (s *ServiceTagWithSemVer) ToGitTag() GitTag {
	return GitTag(s.String())
}
func (s *ServiceTagWithSemVer) String() string {
	return fmt.Sprintf("%s-%s", s.service, s.version.String())
}

func (s *ServiceTagWithSemVer) GreaterThan(other *ServiceTagWithSemVer) bool {
	return s.version.GreaterThan(other.version)
}
func (s *ServiceTagWithSemVer) LessThan(other *ServiceTagWithSemVer) bool {
	return s.version.LessThan(other.version)
}
func (s *ServiceTagWithSemVer) Equal(other *ServiceTagWithSemVer) bool {
	return s.version.Equal(other.version)
}

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

func (s SemVer) GreaterThan(other SemVer) bool {
	if s.major > other.major {
		return true
	}
	if s.major == other.major && s.minor > other.minor {
		return true
	}
	if s.major == other.major && s.minor == other.minor && s.patch > other.patch {
		return true
	}
	return false
}

func (s SemVer) LessThan(other SemVer) bool {
	if s.major < other.major {
		return true
	}
	if s.major == other.major && s.minor < other.minor {
		return true
	}
	if s.major == other.major && s.minor == other.minor && s.patch < other.patch {
		return true
	}
	return false
}

func (s SemVer) Equal(other SemVer) bool {
	return s.major == other.major && s.minor == other.minor && s.patch == other.patch
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
