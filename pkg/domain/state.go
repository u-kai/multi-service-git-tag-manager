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
