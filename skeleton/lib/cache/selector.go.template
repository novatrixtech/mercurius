package cache

import "github.com/go-macaron/cache"

// Option cache selector
func Option(typ string) cache.Options {
	var opt cache.Options
	switch typ {
	case "file":
		opt = File
	case "redis":
		opt = Redis
	case "memcache":
		opt = Memcache
	default:
		opt = Memory
	}
	return opt
}
