package frameworkdto

type Environment int

const (
	EnvProd Environment = iota
	EnvDev
)

type DatabaseType int

const (
	DatabaseTypeMySQL DatabaseType = iota
	DatabaseTypePostgreSQL
	DatabaseTypeSQLite
)

type FrameworkConfig struct {
	Environment Environment    `json:"environment"`
	JWTSecret   string         `json:"jwt_secret"`
	DBType      DatabaseType   `json:"db_type"`
	DbCfg       DatabaseConfig `json:"db_cfg"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"ssl_mode"`
}
