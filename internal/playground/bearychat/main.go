package main

import (
	"time"

	"github.com/jujili/jili/internal/pkg/beary"
)

func main() {
	c := beary.NewChannel()

	c.Verbose("Verbose")
	time.Sleep(time.Millisecond * 321)
	c.Debug("Debug")
	time.Sleep(time.Millisecond * 321)
	c.Info("Info")
	time.Sleep(time.Millisecond * 321)
	c.Warning("Warning")
	time.Sleep(time.Millisecond * 321)
	c.Fatal("Fatal")
}
