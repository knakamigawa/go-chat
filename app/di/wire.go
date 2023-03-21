//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
)

func InitService() (serviceRegistry, error) {
	wire.Build(
		providerDbConnRegistry,
		provideRepositoryRegistry,
		provideServiceRegistry,
	)
	return serviceRegistry{}, nil
}
