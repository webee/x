package app

// Config app configs.
type Config struct {
	Static      bool   `default:"false"`
	DocFile     string `default:"./docs/swagger.json"`
	DocPath     string `default:"/docs/swagger.json"`
	SwaggerPath string `default:"/swagger/"`
	Host        string `default:"localhost:5000"`
}
