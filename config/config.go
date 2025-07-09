package config

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func GetConfig() *Config {
	return &Config{
		DBHost:     "localhost",
		DBPort:     "3306",
		DBUser:     "root",
		DBPassword: "123456",
		DBName:     "user_management",
		ServerPort: "8080",
	}
}
