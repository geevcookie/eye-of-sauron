package internal

// Config represents the application config
type Config struct {
	ServerConfig  ServerConfig  `mapstructure:"server"`
	MetricsConfig MetricsConfig `mapstructure:"metrics"`
}

// ServerConfig represents the API server config
type ServerConfig struct {
	Port string
}

// MetricsConfig represents the Metrics Scraper config
type MetricsConfig struct {
	Frequency int
}
