package injection

import (
	"github.com/hcsouza/bard/cache"
	. "github.com/hcsouza/bard/logger"
	"github.com/sarulabs/di"
)

var Services = []di.Def{
	{
		Name:  "CacheClient",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return cache.NewCacheClient(), nil
		},
	}}

var (
	appContainer di.Container
)

func init() {
	appContainer = CreateContainer()
}

func SetContainerApp(container di.Container) {
	appContainer = container
}

func CreateContainer() di.Container {
	builder, err := di.NewBuilder()
	if err != nil {
		Logger.Error(err.Error())
	}
	_ = builder.Add(Services...)
	return builder.Build()
}

func Get(dep string) interface{} {
	return di.Get(appContainer, dep)
}
