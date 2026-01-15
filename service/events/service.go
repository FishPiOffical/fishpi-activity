package events

import (
	"log/slog"

	"github.com/pocketbase/pocketbase/core"
)

type Service struct {
	app core.App

	logger *slog.Logger
}

func NewService(app core.App) *Service {
	service := &Service{
		app:    app,
		logger: app.Logger().WithGroup("service.events"),
	}
	return service
}
