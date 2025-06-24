package main

import (
	"time"

	"github.com/scottEAdams1/REPLPokedex2/internal/pokecache"
)

func main() {
	cache := pokecache.NewCache(time.Minute * 5)
	startREPL(cache)
}
