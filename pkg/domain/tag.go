package domain

import (
	"fmt"
	"regexp"
)

type GitTag string

func (g GitTag) String() string {
	return string(g)
}

type ServiceTagWithSemVer struct {
	Service ServiceName
	Version SemVer
}

func NewServiceTagWithSemVer(service ServiceName, version SemVer) *ServiceTagWithSemVer {
	return &ServiceTagWithSemVer{
		Service: service,
		Version: version,
	}
}

func (s *ServiceTagWithSemVer) UpdateMajor() {
	s.Version = s.Version.MajorUp()
}
func (s *ServiceTagWithSemVer) UpdateMinor() {
	s.Version = s.Version.MinorUp()
}
func (s *ServiceTagWithSemVer) UpdatePatch() {
	s.Version = s.Version.PatchUp()
}
func (s *ServiceTagWithSemVer) ToGitTag() GitTag {
	return GitTag(s.String())
}
func (s *ServiceTagWithSemVer) String() string {
	return fmt.Sprintf("%s-%s", s.Service, s.Version.String())
}

func (s *ServiceTagWithSemVer) GreaterThan(other *ServiceTagWithSemVer) bool {
	return s.Version.GreaterThan(other.Version)
}
func (s *ServiceTagWithSemVer) LessThan(other *ServiceTagWithSemVer) bool {
	return s.Version.LessThan(other.Version)
}
func (s *ServiceTagWithSemVer) Equal(other *ServiceTagWithSemVer) bool {
	return s.Version.Equal(other.Version)
}

type SemVer struct {
	Major int
	Minor int
	Patch int
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
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func (s SemVer) GreaterThan(other SemVer) bool {
	if s.Major > other.Major {
		return true
	}
	if s.Major == other.Major && s.Minor > other.Minor {
		return true
	}
	if s.Major == other.Major && s.Minor == other.Minor && s.Patch > other.Patch {
		return true
	}
	return false
}

func (s SemVer) LessThan(other SemVer) bool {
	if s.Major < other.Major {
		return true
	}
	if s.Major == other.Major && s.Minor < other.Minor {
		return true
	}
	if s.Major == other.Major && s.Minor == other.Minor && s.Patch < other.Patch {
		return true
	}
	return false
}

func (s SemVer) Equal(other SemVer) bool {
	return s.Major == other.Major && s.Minor == other.Minor && s.Patch == other.Patch
}

func (s SemVer) MajorUp() SemVer {
	return SemVer{
		Major: s.Major + 1,
		Minor: 0,
		Patch: 0,
	}
}

func (s SemVer) MinorUp() SemVer {
	return SemVer{
		Major: s.Major,
		Minor: s.Minor + 1,
		Patch: 0,
	}
}

func (s SemVer) PatchUp() SemVer {
	return SemVer{
		Major: s.Major,
		Minor: s.Minor,
		Patch: s.Patch + 1,
	}
}

func (s *SemVer) String() string {
	return fmt.Sprintf("v%d.%d.%d", s.Major, s.Minor, s.Patch)
}

var (
	serviceSemVerRe         = regexp.MustCompile(`^([a-zA-Z0-9-]+)-v(\d+)\.(\d+)\.(\d+)$`)
	serviceSemVerWithoutVRe = regexp.MustCompile(`^([a-zA-Z0-9-]+)-(\d+)\.(\d+)\.(\d+)$`)
)

func errInvalidServiceSemVerMsg(invalid string) error {
	return fmt.Errorf("invalid service semver string: %s\nservice Version should be SERVICE_NAME-vMAJOR.MINOR.PATCH", invalid)
}

func (g GitTag) ToServiceTag() (*ServiceTagWithSemVer, error) {
	matches := serviceSemVerRe.FindStringSubmatch(g.String())
	semVerStr := "v"
	if len(matches) != 5 {
		matches = serviceSemVerWithoutVRe.FindStringSubmatch(g.String())
		if len(matches) != 5 {
			return nil, errInvalidServiceSemVerMsg(g.String())
		}
		// without v Version
		semVerStr = ""
	}

	service := matches[1]
	Version, err := FromStr(fmt.Sprintf("%s%s.%s.%s", semVerStr, matches[2], matches[3], matches[4]))
	if err != nil {
		return nil, errInvalidServiceSemVerMsg(g.String())
	}

	return NewServiceTagWithSemVer(ServiceName(service), Version), nil
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

type VersionUpServiceTag func(*[]GitTag) *[]*ServiceTagWithSemVer

func MajorUpAll(tags *[]GitTag) *[]*ServiceTagWithSemVer {
	return VersionUpAll(func(tag *ServiceTagWithSemVer) {
		tag.UpdateMajor()
	})(tags)
}

func MinorUpAll(tags *[]GitTag) *[]*ServiceTagWithSemVer {
	return VersionUpAll(func(tag *ServiceTagWithSemVer) {
		tag.UpdateMinor()
	})(tags)
}

func PatchUpAll(tags *[]GitTag) *[]*ServiceTagWithSemVer {
	return VersionUpAll(func(tag *ServiceTagWithSemVer) {
		tag.UpdatePatch()
	})(tags)
}

type VersionUpFunc func(*ServiceTagWithSemVer)

func VersionUpAll(f VersionUpFunc) VersionUpServiceTag {
	return func(tags *[]GitTag) *[]*ServiceTagWithSemVer {
		serviceTags := []*ServiceTagWithSemVer{}
		if tags == nil || len(*tags) == 0 {
			return &serviceTags
		}

		tmpAlreadyUpdatedServiceTags := map[ServiceName]*ServiceTagWithSemVer{}

		for _, tag := range *tags {
			serviceTag, err := tag.ToServiceTag()
			if err != nil {
				continue
			}
			f(serviceTag)
			if Version, ok := tmpAlreadyUpdatedServiceTags[serviceTag.Service]; ok {
				if serviceTag.LessThan(Version) || serviceTag.Equal(Version) {
					continue
				}
				if serviceTag.GreaterThan(Version) {
					tmpAlreadyUpdatedServiceTags[serviceTag.Service] = serviceTag
					continue
				}
			}
			tmpAlreadyUpdatedServiceTags[serviceTag.Service] = serviceTag
		}

		for _, serviceTag := range tmpAlreadyUpdatedServiceTags {
			serviceTags = append(serviceTags, serviceTag)
		}

		return &serviceTags
	}
}

func SortsServiceTags(tags *[]*ServiceTagWithSemVer) map[ServiceName][]*ServiceTagWithSemVer {
	sorted := map[ServiceName][]*ServiceTagWithSemVer{}
	for _, tag := range *tags {
		if _, ok := sorted[tag.Service]; !ok {
			sorted[tag.Service] = []*ServiceTagWithSemVer{}
		}
		sorted[tag.Service] = append(sorted[tag.Service], tag)
	}
	for _, serviceTags := range sorted {
		for i := 0; i < len(serviceTags); i++ {
			for j := i + 1; j < len(serviceTags); j++ {
				if serviceTags[i].GreaterThan(serviceTags[j]) {
					serviceTags[i], serviceTags[j] = serviceTags[j], serviceTags[i]
				}
			}
		}
	}
	return sorted
}
