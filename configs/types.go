package configs

type Config struct {
	// Server
	ServerPort string `mapstructure:"SERVER_PORT"` 
	ServerEnv     string `mapstructure:"SERVER_ENV"`  

	// Database
	DBHost string `mapstructure:"DB_HOST"`    
	DBPort string `mapstructure:"DB_PORT"`    
	DBUser string `mapstructure:"DB_USER"`    
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName string `mapstructure:"DB_NAME"`    
	DBSSLMode string `mapstructure:"DB_SSLMODE"`

	// JWT
	JWTSecret        string `mapstructure:"JWT_SECRET"` 
	JWTDuration int    `mapstructure:"JWT_DURATION"` 

	// Redis
	RedisHost string `mapstructure:"REDIS_HOST"` 
	RedisPort string `mapstructure:"REDIS_PORT"` 

	// Log
	LogLevel string `mapstructure:"LOG_LEVEL"` 
	LogFormat string `mapstructure:"LOG_FORMAT"` 

	// Stripe
	StripeAPIKey      string `mapstructure:"STRIPE_API_KEY"`      
	StripeWebhookSecret string `mapstructure:"STRIPE_WEBHOOK_SECRET"`
}
