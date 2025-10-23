package application

import (
	"github.com/pocketbase/pocketbase/core"
)

type fixBugHandler func(e *core.BootstrapEvent) error

func (application *Application) fixBug(e *core.BootstrapEvent) error {
	list := []fixBugHandler{
		application.fixExample,
	}

	for _, handler := range list {
		if err := handler(e); err != nil {
			return err
		}
	}

	return nil
}

func (application *Application) fixExample(*core.BootstrapEvent) error {
	return nil
}
