package usecase

import (
	"msgtm/pkg/domain"
)

// In the case of using APIs that involve side effects such as registration or deletion.
type CommandExecutor[Command any] interface {
	Execute(Command) error
}

// In the case of using APIs that do not involve side effects such as queries.
type QueryExecutor[Query any, Entity any] interface {
	Execute(Query) (Entity, error)
}

/// Executors and Comands

// DestroyServiceTags is a usecase that deletes the specified tags.
type DestroyServiceTags = CommandExecutor[DestroyServiceTagsCommand]
type DestroyServiceTagsCommand struct {
	Tags *[]*domain.ServiceTagWithSemVer
}

// CommitGetter is a usecase that gets the tags of the specified commit.
type ListTags = QueryExecutor[ListTagsQuery, *[]domain.GitTag]
type ListTagsQuery struct {
	Filter func(*domain.ServiceName) bool
}

// CommitGetter is a usecase that gets the tags of the specified commit.
type CommitTagGetter = QueryExecutor[GetCommitTagQuery, *[]domain.GitTag]
type GetCommitTagQuery struct {
	CommitId *domain.CommitId
}

// CommitGetter is a usecase that gets the tags of the specified commit.
type CommitFinder = QueryExecutor[FindCommitQuery, *domain.CommitId]
type FindCommitQuery struct {
	Tag *domain.GitTag
}

// CommitPusher is a usecase that pushes the specified tags to the remote repository.
type CommitPusher = CommandExecutor[CommitPushCommand]
type CommitPushCommand struct {
	RemoteAddr *domain.RemoteAddr
	Tags       *[]*domain.ServiceTagWithSemVer
}

// RegisterServiceTags is a usecase that registers the specified tags.
type RegisterServiceTags = CommandExecutor[RegisterServiceTagsCommand]
type RegisterServiceTagsCommand struct {
	CommitId *domain.CommitId
	Tags     *[]*domain.ServiceTagWithSemVer
}
