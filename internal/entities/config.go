package entities

type Config struct {
	User        string `env:"POSTGRES_USER"`
	Password    string `env:"POSTGRES_PASSWORD"`
	Host        string `env:"DB_HOST"`
	Port        string `env:"DB_PORT"`
	DBName      string `env:"POSTGRES_DB"`
	MaxAttempts int    `env:"MAX_ATTEMPTS"`
}
