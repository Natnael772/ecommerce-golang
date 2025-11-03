package configs

const (
    // Server
    ServerPort = "SERVER_PORT"
    ServerEnv     = "SERVER_ENV"

    // Database
    DBHost     = "DB_HOST"
    DBPort     = "DB_PORT"
    DBUser     = "DB_USER"
    DBPassword = "DB_PASSWORD"
    DBName     = "DB_NAME"
    DBSSLMode  = "DB_SSLMODE"

    // Auth
    JWTSecret     = "JWT_SECRET"
    JWTDuration   = "JWT_DURATION"

    // Cache
    RedisHost = "REDIS_HOST"
    CacheRedisPort = "REDIS_PORT"

    // Log
    LogLevel  = "log.level"
    LogFormat = "log.format"

    // Stripe
    StripeAPIKey       = "STRIPE_API_KEY"
    StripeWebhookSecret = "STRIPE_WEBHOOK_SECRET"
)
