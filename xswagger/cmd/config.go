package cmd

import (
	"github.com/webee/x/xswagger/app"
)

// Config is the all settings of this command
type Config struct {
	Debug   bool   `default:"false"`
	Address string `default:":7070"`

	APP app.Config
}
