package sql

type DatabaseConfigType struct {
	Host     string `json:"DB_HOST"`
	Port     string `json:"DB_PORT"`
	UserName string `json:"DB_USER"`
	Password string `json:"DB_PASSWORD"`
	Database string `json:"DB_NAME"`
	Protocol string `json:"protocol"`
	Timeout  int    `json:"timeout"`
	MaxOpen  int    `json:"max_open"`
	MaxIdle  int    `json:"max_idle"`
	LogFile  string `json:"log_file"`
	DBType   string `yaml:"dbtype"`
}
