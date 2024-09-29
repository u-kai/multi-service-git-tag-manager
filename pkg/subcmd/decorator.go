package subcmd

import "log/slog"

type SubCommand[T any] func(T) error

func LogSubCommandDecorator[T any](f SubCommand[T], logger *slog.Logger) SubCommand[T] {
	return func(cmd T) error {
		logger.Debug("SubCommand: Execute", slog.Any("param", cmd))
		err := f(cmd)
		if err != nil {
			logger.Error("SubCommand: Execute", slog.Any("error", err))
			return err
		}
		return nil
	}
}
