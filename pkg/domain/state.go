package domain

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

type ServiceTagState struct {
	ServiceName *ServiceName    `json:"name" yaml:"name"`
	Latest      *ServiceTagInfo `json:"latest" yaml:"latest"`
	Prev        *ServiceTagInfo `json:"prev" yaml:"prev"`
}

func InitServiceTagState(serviceName *ServiceName) *ServiceTagState {
	return &ServiceTagState{
		ServiceName: serviceName,
		Latest:      nil,
		Prev:        nil,
	}
}

func (state *ServiceTagState) UpdateLatest(tag *ServiceTagInfo) {
	state.Prev = state.Latest
	state.Latest = tag
}

type ServiceTagInfo struct {
	Tag           *ServiceTagWithSemVer `json:"tag" yaml:"tag"`
	CommitId      *CommitId             `json:"commitId" yaml:"commitId"`
	Description   *string               `json:"description" yaml:"description"`
	CommitComment *string               `json:"commitComment" yaml:"commitComment"`
}

type WritedState struct {
	ServiceTagStates []*ServiceTagState `json:"services" yaml:"services"`
}

func InitStateWriter(services ...ServiceName) *WritedState {
	states := make([]*ServiceTagState, 0)
	for _, service := range services {
		states = append(states, InitServiceTagState(&service))
	}
	return &WritedState{
		ServiceTagStates: states,
	}
}

type WriteFormat int

const (
	JSON WriteFormat = iota
	YAML
)

func (s *WritedState) Update(serviceName ServiceName, tag *ServiceTagInfo) {
	// len 0はfor文が実行されないため
	if len(s.ServiceTagStates) == 0 {
		s.ServiceTagStates = append(s.ServiceTagStates, InitServiceTagState(&serviceName))
		s.ServiceTagStates[0].UpdateLatest(tag)
		return
	}
	for i, state := range s.ServiceTagStates {
		if *state.ServiceName == serviceName {
			state.UpdateLatest(tag)
			return
		}
		if i == len(s.ServiceTagStates)-1 {
			newS := InitServiceTagState(&serviceName)
			newS.UpdateLatest(tag)
			s.ServiceTagStates = append(s.ServiceTagStates, newS)
		}
	}
}

func (s *WritedState) Write(writer io.Writer, format WriteFormat) error {
	switch format {
	case JSON:
		b, err := json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = writer.Write(b)
		if err != nil {
			return fmt.Errorf("failed to write state: %w", err)
		}
		return nil
	case YAML:
		b, err := yaml.Marshal(s)
		if err != nil {
			return err
		}
		_, err = writer.Write(b)
		if err != nil {
			return fmt.Errorf("failed to write state: %w", err)
		}
	default:
		return fmt.Errorf("unknown format: %d", format)
	}

	return nil
}

func FromReader(reader io.Reader, format WriteFormat) (*WritedState, error) {
	b := []byte{}
	_, err := reader.Read(b)
	if err != nil {
		return nil, err
	}
	switch format {
	case JSON:
		state := &WritedState{}
		err = json.Unmarshal(b, state)
		if err != nil {
			return nil, err
		}
		return state, nil
	case YAML:
		state := &WritedState{}
		err = yaml.Unmarshal(b, state)
		if err != nil {
			return nil, err
		}
		return state, nil
	}
	return nil, fmt.Errorf("unsupported format: %v", format)
}

type marshaledState struct {
	Services []struct {
		Name   string                    `json:"name" yaml:"name"`
		Latest *marshaledServiceTagState `json:"latest" yaml:"latest"`
		Prev   *marshaledServiceTagState `json:"prev" yaml:"prev"`
	} `json:"services" yaml:"services"`
}

type marshaledServiceTagState struct {
	Tag struct {
		Version string `json:"version" yaml:"version"`
	} `json:"tag" yaml:"tag"`
	CommitId      string `json:"commitId" yaml:"commitId"`
	Description   string `json:"description" yaml:"description"`
	CommitComment string `json:"commitComment" yaml:"commitComment"`
}

func (s *WritedState) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.toMarshaled())
}
func (s *WritedState) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(s.toMarshaled())
}

func (s *WritedState) UnmarshalJSON(b []byte) error {
	m := marshaledState{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	return s.fromMarshaled(m)

}

func (s *WritedState) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tmp := marshaledState{}
	if err := unmarshal(&tmp); err != nil {
		return err
	}
	return s.fromMarshaled(tmp)
}

func (s *WritedState) fromMarshaled(m marshaledState) error {
	states := make([]*ServiceTagState, 0, len(m.Services))
	for _, service := range m.Services {
		name := ServiceName(service.Name)
		if service.Latest == nil && service.Prev == nil {
			states = append(states, InitServiceTagState(&name))
			continue
		}
		if service.Latest != nil {
			version, err := FromStr(service.Latest.Tag.Version)
			if err != nil {
				return err
			}
			serviceTag := NewServiceTagWithSemVer(name, version)
			commitId := CommitId(service.Latest.CommitId)
			var description *string = nil
			if service.Latest.Description != "" {
				description = &service.Latest.Description
			}
			var commitComment *string = nil
			if service.Latest.CommitComment != "" {
				commitComment = &service.Latest.CommitComment
			}
			state := InitServiceTagState(&name)
			state.UpdateLatest(
				&ServiceTagInfo{
					Tag:           serviceTag,
					CommitId:      &commitId,
					Description:   description,
					CommitComment: commitComment,
				},
			)
			states = append(states, state)
		}
		if service.Prev != nil {
			version, err := FromStr(service.Prev.Tag.Version)
			if err != nil {
				return err
			}
			serviceTag := NewServiceTagWithSemVer(name, version)
			commitId := CommitId(service.Prev.CommitId)
			state := InitServiceTagState(&name)
			var description *string = nil
			if service.Prev.Description != "" {
				description = &service.Prev.Description
			}
			var commitComment *string = nil
			if service.Prev.CommitComment != "" {
				commitComment = &service.Prev.CommitComment
			}
			state.Prev = &ServiceTagInfo{
				Tag:           serviceTag,
				CommitId:      &commitId,
				Description:   description,
				CommitComment: commitComment,
			}
			states = append(states, state)
		}
	}
	s.ServiceTagStates = states
	return nil
}

func (s *WritedState) toMarshaled() marshaledState {
	services := make([]struct {
		Name   string                    `json:"name" yaml:"name"`
		Latest *marshaledServiceTagState `json:"latest" yaml:"latest"`
		Prev   *marshaledServiceTagState `json:"prev" yaml:"prev"`
	}, len(s.ServiceTagStates))
	m := marshaledState{
		Services: services,
	}
	for i, state := range s.ServiceTagStates {
		m.Services[i].Name = state.ServiceName.String()
		if state.Latest != nil {
			m.Services[i].Latest = &marshaledServiceTagState{
				Tag: struct {
					Version string `json:"version" yaml:"version"`
				}{
					Version: state.Latest.Tag.String(),
				},
				CommitId:      state.Latest.CommitId.String(),
				Description:   *state.Latest.Description,
				CommitComment: *state.Latest.CommitComment,
			}
		}
		if state.Prev != nil {
			m.Services[i].Prev = &marshaledServiceTagState{
				Tag: struct {
					Version string `json:"version" yaml:"version"`
				}{
					Version: state.Prev.Tag.String(),
				},
				CommitId:      state.Prev.CommitId.String(),
				Description:   *state.Prev.Description,
				CommitComment: *state.Prev.CommitComment,
			}
		}
	}
	return m
}
