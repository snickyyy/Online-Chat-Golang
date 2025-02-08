//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"libs/src/settings"
)

func InitializeApp() (*settings.App, error) {
	wire.Build(
		settings.GetBaseConfig,
		settings.InitDb,
		settings.InitLogger,
		settings.NewApp,
	)
	return nil, nil
}