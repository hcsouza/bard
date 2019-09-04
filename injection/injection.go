package injection

import (
	"github.com/hcsouza/bard/cache"
	"github.com/sarulabs/di"
	"log"
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

func CreateContainer() di.Container {
	builder, err := di.NewBuilder()
	if err != nil {
		log.Println(err.Error())
	}
	_ = builder.Add(Services...)
	return builder.Build()
}

func Get(dep string) interface{} {
	return di.Get(appContainer, dep)
}
