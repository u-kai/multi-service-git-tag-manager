package usecase_test

import (
	"msgtm/pkg/domain"
	"msgtm/pkg/usecase"
	"reflect"
	"testing"
)

func TestCreateServiceTags(t *testing.T) {
	services := []domain.ServiceName{"service-a", "service-b"}
	version := domain.NewSemVer(0, 0, 1)
	commitId := domain.HEAD
	mockRegister := &MockRegister{}
	err := usecase.CreateServiceTags(mockRegister, &commitId, services, version)
	if err != nil {
		t.Errorf("CreateServiceTags() error = %v, want nil", err)
	}
	expected := []*domain.ServiceTagWithSemVer{
		domain.NewServiceTagWithSemVer("service-a", domain.NewSemVer(0, 0, 1)),
		domain.NewServiceTagWithSemVer("service-b", domain.NewSemVer(0, 0, 1)),
	}
	if !cmpArrayContent(
		*mockRegister.AddedTags,
		expected,
	) {
		t.Errorf(
			"CreateServiceTags() = %v, want %v", mockRegister.AddedTags, expected,
		)
	}
}

// 順不同な配列の比較
func cmpArrayContent[T any](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		found := false
		for _, vv := range b {
			if reflect.DeepEqual(v, vv) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
