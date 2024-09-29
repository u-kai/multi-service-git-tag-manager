package usecase

import (
	"log/slog"
	"msgtm/pkg/domain"
)

type LoggingCommandExecutor[Command any] struct {
	Executor CommandExecutor[Command]
	Logger   slog.Logger
}

func (l *LoggingCommandExecutor[Command]) Execute(cmd Command) error {
	l.Logger.Debug("CommandExecutor: Execute", slog.Any("command", cmd))
	err := l.Executor.Execute(cmd)
	if err != nil {
		l.Logger.Error("CommandExecutor: Execute", slog.Any("error", err))
		return err
	}
	l.Logger.Debug("CommandExecutor: Execute", slog.Any("result", "success"))
	return nil
}

type LoggingQueryExecutor[Query any, Entity any] struct {
	Executor QueryExecutor[Query, Entity]
	Logger   slog.Logger
}

func (l *LoggingQueryExecutor[Query, Entity]) Execute(query Query) (Entity, error) {
	l.Logger.Debug("QueryExecutor: Execute", slog.Any("query", query))
	entity, err := l.Executor.Execute(query)
	if err != nil {
		l.Logger.Error("QueryExecutor: Execute", slog.Any("error", err))
		e := new(Entity)
		return *e, err
	}
	l.Logger.Debug("QueryExecutor: Execute", slog.Any("result", entity))
	return entity, nil
}

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
	Filter *[]*domain.ServiceName
}

// CommitGetter is a usecase that gets the tags of the specified commit.
type CommitTagGetter = QueryExecutor[GetCommitTagQuery, *[]domain.GitTag]
type GetCommitTagQuery struct {
	CommitId *domain.CommitId
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
