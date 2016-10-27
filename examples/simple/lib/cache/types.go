package cache

import (
	"github.com/go-macaron/cache"
	_ "github.com/go-macaron/cache/memcache"
	_ "github.com/go-macaron/cache/redis"
	"github.com/novatrixtech/mercurius/examples/simple/conf"
)

var (
	Memory = cache.Options{}
	File   = cache.Options{
		Adapter:       "file",
		AdapterConfig: conf.Cfg.Section("").Key("cache_adpter_config").Value(),
	}
	Redis = cache.Options{
		Adapter:       "redis",
		AdapterConfig: conf.Cfg.Section("").Key("cache_adpter_config").Value(),
	}
	Memcache = cache.Options{
		Adapter:       "memcache",
		AdapterConfig: conf.Cfg.Section("").Key("cache_adpter_config").Value(),
	}
)
