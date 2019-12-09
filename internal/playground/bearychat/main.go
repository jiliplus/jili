package main

import "github.com/aQuaYi/jili/internal/pkg/beary"

import "time"

func main() {
	c := beary.NewChannel()

	c.Verbose("verbose")
	time.Sleep(time.Millisecond * 321)
	c.Debug("Debug")
	time.Sleep(time.Millisecond * 321)
	c.Info("Info")
	time.Sleep(time.Millisecond * 321)
	c.Warning("Warning")
	time.Sleep(time.Millisecond * 321)
	c.Fatal("Fatal")
}
