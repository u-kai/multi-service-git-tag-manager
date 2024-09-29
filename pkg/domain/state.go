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

type StateWriter struct {
	ServiceTagStates []*ServiceTagState
}

func InitStateWriter(services ...ServiceName) *StateWriter {
	states := make([]*ServiceTagState, 0)
	for _, service := range services {
		states = append(states, InitServiceTagState(&service))
	}
	return &StateWriter{
		ServiceTagStates: states,
	}
}

type WriteFormat int

const (
	JSON WriteFormat = iota
	YAML
)

type WritedState struct {
	Services []*ServiceTagState `json:"services" yaml:"services"`
}

func (stateWriter *StateWriter) Write(writer io.Writer, format WriteFormat) error {
	switch format {
	case JSON:
		b, err := json.Marshal(WritedState{Services: stateWriter.ServiceTagStates})
		if err != nil {
			return err
		}
		_, err = writer.Write(b)
		if err != nil {
			return fmt.Errorf("failed to write state: %w", err)
		}
		return nil
	case YAML:
		b, err := yaml.Marshal(WritedState{Services: stateWriter.ServiceTagStates})
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
