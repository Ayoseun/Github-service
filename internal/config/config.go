package config

type DatabaseConfig struct {
	DSN string
}

func GetDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		DSN: "postgres://ayoseun:Jared15$@localhost:5432/github",
	}
}
