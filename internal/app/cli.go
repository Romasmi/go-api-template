package app

import (
	"github.com/Romasmi/s-shop-microservices/internal/interface/cli"
)

type Cli struct {
	*App
}

func NewCli(app *App) *Cli {
	return &Cli{App: app}
}

func (c *Cli) Run() error {
	cli.SetApp(c.App)
	return cli.Execute()
}
