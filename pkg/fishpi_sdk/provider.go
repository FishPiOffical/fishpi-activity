package fishpi_sdk

import (
	"bless-activity/model"
	"encoding/json"

	"github.com/FishPiOffical/golang-sdk/config"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Provider struct {
	app core.App

	record *model.Config
	config *config.Config
}

func NewProvider(app core.App) (*Provider, error) {
	record := new(model.Config)
	if err := app.RecordQuery(model.DbNameConfigs).Where(dbx.HashExp{
		model.ConfigsFieldKey: model.ConfigKeyFishpi,
	}).One(record); err != nil {
		return nil, err
	}

	cfg := new(config.Config)
	if err := json.Unmarshal([]byte(record.Value()), cfg); err != nil {
		return nil, err
	}

	provider := &Provider{
		app: app,

		record: record,
		config: cfg,
	}

	return provider, nil
}

func (provider *Provider) Get() *config.Config {
	return provider.config
}

func (provider *Provider) Update(cfg *config.Config) error {
	provider.config = cfg

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	provider.record.SetValue(string(data))
	if err = provider.app.Save(provider.record); err != nil {
		return err
	}

	return nil
}
