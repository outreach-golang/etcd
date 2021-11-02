package etcd

type Option func(*Config)

type Config struct {
	//环境变量
	env string
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
		env:           "default",
		points:        []string{"127.0.0.1:2379"},
		needSSL:       1,
		sslServerFile: "",
		sslKeyFile:    "",
		caFile:        "",
	}
}

func Env(env string) Option {
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
