package etcd

type EnvVar string

const (
	DefaultEV    EnvVar = "default"
	TestingEV    EnvVar = "testing"
	ReleaseEV    EnvVar = "release"
	ProductionEV EnvVar = "production"
)

type Option func(*Config)

type Config struct {
	//环境变量
	env EnvVar
	//节点
	points []string
	//是否需要ssl
	needSSL       int
	dirPath       string
	sslServerFile string
	sslKeyFile    string
	caFile        string
}

func DefaultConf() *Config {
	return &Config{
		env:           DefaultEV,
		points:        []string{"127.0.0.1:2379"},
		needSSL:       1,
		dirPath:       "./configs/default/",
		sslServerFile: "server.pem",
		sslKeyFile:    "server-key.pem",
		caFile:        "ca.pem",
	}
}

func Env(env EnvVar) Option {
	return func(config *Config) {
		config.env = env
	}
}

func DirPath(dirPath string) Option {
	return func(config *Config) {
		config.dirPath = dirPath
	}
}

func Points(p []string) Option {
	return func(config *Config) {
		config.points = p
	}
}

func NeedSSL(n int) Option {
	return func(config *Config) {
		config.needSSL = n
	}
}

func SSLServerFile(filePath string) Option {
	return func(config *Config) {
		config.sslServerFile = filePath
	}
}

func SSLKeyFile(filePath string) Option {
	return func(config *Config) {
		config.sslKeyFile = filePath
	}
}

func SSLCaFile(filePath string) Option {
	return func(config *Config) {
		config.caFile = filePath
	}
}
