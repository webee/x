package cmd

import (
	"flag"

	"github.com/webee/x/xconfig"
	"github.com/webee/x/xswagger/app"
)

var (
	config = new(Config)
	a      *app.App
)

func init() {
	rv := new(xconfig.SyncValues)
	xconfig.Load(config, func() {
		flagConfig := new(Config)
		flag.BoolVar(rv.AddBool(&flagConfig.Debug, &config.Debug), "debug", false, "debug?")
		flag.StringVar(rv.AddString(&flagConfig.Address, &config.Address), "address", ":7070", "listening address.")
		flag.BoolVar(rv.AddBool(&flagConfig.APP.Static, &config.APP.Static), "s", false, "static mode, do not need swagger.json")
		flag.StringVar(rv.AddString(&flagConfig.APP.DocFile, &config.APP.DocFile), "doc", "./docs/swagger.json", "swagger.json file path.")
		flag.StringVar(rv.AddString(&flagConfig.APP.Host, &config.APP.Host), "host", "localhost:5000", "api host.")
	})
	// sync make flags override config file.
	rv.Sync()

	// app
	a = app.CreateApp(config.Debug, &config.APP)
}

// GetConfig 返回当前配置
func GetConfig() *Config {
	return config
}

// Start server
func Start() {
	// change to real host when use 'doc.json'.
	// Note: do it in your command
	// docs.SwaggerInfo.Host = cmd.GetConfig().APP.Host

	defer a.Destroy()
	a.Run(config.Address)
}
