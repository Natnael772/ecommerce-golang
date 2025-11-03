package configs

import (
	"strings"

	"github.com/spf13/viper"
)

func BindAllKeys() {
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    keys := []string{
        // Server
        ServerPort,
        ServerEnv,

		// Database
		DBHost,
		DBPort,
		DBUser,
		DBPassword,
		DBName,
		DBSSLMode,

		// Auth
		JWTSecret,
		JWTDuration,

		// Cache
		RedisHost,
		CacheRedisPort,

		// Log
		LogLevel,
		LogFormat,

		// Stripe
		StripeAPIKey,
		StripeWebhookSecret,
    }

    for _, key := range keys {
        _ = viper.BindEnv(key)
    }
}
