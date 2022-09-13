package cmd

import (
	"flag"

	"github.com/swaggo/swag"
	"github.com/webee/x/xconfig"
	"github.com/webee/x/xswagger/app"
)

var (
	config = new(Config)
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

}

// GetConfig 返回当前配置
func GetConfig() *Config {
	return config
}

// Start server
func Start(swaggerInfo *swag.Spec) {
	// app
	app := app.CreateApp(config.Debug, &config.APP, swaggerInfo)
	defer app.Destroy()

	app.Run(config.Address)
}
