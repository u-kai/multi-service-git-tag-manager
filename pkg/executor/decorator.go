package executor

import (
	"encoding/json"
	"log/slog"
	"msgtm/pkg/usecase"
)

type LoggingCommandExecutor[Command any] struct {
	Executor usecase.CommandExecutor[Command]
	Logger   *slog.Logger
}

type stringer interface {
	String() string
}

func (l *LoggingCommandExecutor[Command]) Execute(cmd Command) error {
	var cmdAny any = cmd
	switch cmdAny := cmdAny.(type) {
	case stringer:
		l.Logger.Debug("UsecaseCommandExecutor: Execute", slog.String("command", cmdAny.String()))
	default:
		cmdJson, _ := json.Marshal(cmd)
		l.Logger.Debug("UsecaseCommandExecutor: Execute", slog.Any("command", string(cmdJson)))
	}
	err := l.Executor.Execute(cmd)
	if err != nil {
		l.Logger.Error("UsecaseCommandExecutor: Execute", slog.Any("error", err))
		return err
	}
	l.Logger.Debug("UsecaseCommandExecutor: Execute", slog.Any("result", "success"))
	return nil
}

type LoggingQueryExecutor[Query any, Entity any] struct {
	Executor usecase.QueryExecutor[Query, Entity]
	Logger   *slog.Logger
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
