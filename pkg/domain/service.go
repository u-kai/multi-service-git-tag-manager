package domain

import "strings"

type ServiceName string

func (s *ServiceName) String() string {
	return string(*s)
}

func (s *ServiceName) IsServiceTag(tag *GitTag) bool {
	return strings.HasPrefix(string(*tag), s.String())
}
