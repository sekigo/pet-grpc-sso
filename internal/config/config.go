package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct { // struct теги метаинформация, которую может использовать любой go код
	Env         string        `yaml:"env" env-default:"local" `
	StoragePath string        `yaml:"storage_path" env-requuired:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:port`
	Timeout time.Duration `yaml:timeout`
}

// соглашение по Must в названии говорит о том, что мы не будем возвращать ошибку
// если например не хочется обрабатывать ошибку.
func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("Config is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config fil does not exist")
	}
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed yo read config: " + err.Error())
	}

	return &cfg

}

// эта функция позволяет использовать задаваемые флагом переменные окружения. Если нет флагов, тогда
// мы вернем то, что прописано в нашем конфиге
func fetchConfigPath() string {
	var res string // сюда будет записано имя флага. Указатель передается, чтобы внутри можно было
	// перезаписать значение res

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse() // выполняем парсинг флаг

	if res == "" { //  если рес пустой, то мы прочитаем переменную приложения
		return os.Getenv("CONFIG_PATH")
	}

	return res
}
