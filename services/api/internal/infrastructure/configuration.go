package infrastructure

type Configuration struct {
	Port              string   `env:"PORT" default:"8080"`
	Address           string   `env:"ADDRESS" default:"localhost"`
	AllowedOrigins    []string `env:"ALLOWED_ORIGINS" separator:","`
	JWTSecret         string   `env:"JWT_SECRET" default:"SECRET"`
	InterserverSecret string   `env:"INTERSERVER_SECRET" default:"SECRET"`

	MongoDatabaseName     string `env:"MONGO_DATABASE" default:"apiDatabase"`
	MongoDatabaseUserName string `env:"MONGO_USER" default:"apiUser"`
	MongoDatabasePassword string `env:"MONGO_PASSWORD" default:"apiUserPassword"`
	MongoDatabaseHost     string `env:"MONGO_HOST" default:"localhost"`
	MongoDatabasePort     string `env:"MONGO_PORT" default:"27017"`
}
