package fishpi

import (
	"bless-activity/model"
	"encoding/json"
	"log/slog"

	"github.com/imroc/req/v3"
	"github.com/lxzan/gws"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Service struct {
	config *Config

	conn *gws.Conn

	app    core.App
	client *req.Client
	logger *slog.Logger
}

func NewService(app core.App) (*Service, error) {
	s := &Service{
		app:    app,
		logger: app.Logger().WithGroup("fishpi"),
	}

	if err := s.init(); err != nil {
		return nil, err
	}

	return s, nil
}

func (service *Service) init() error {
	config := new(model.Config)
	if err := service.app.RecordQuery(model.DbNameConfigs).Where(dbx.HashExp{model.ConfigsFieldKey: model.ConfigKeyFishpi}).One(config); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(config.Value()), &service.config); err != nil {
		return err
	}

	service.client = req.NewClient().
		SetBaseURL(service.config.BaseUrl).
		SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko)")

	//go service.start()

	return nil
}
