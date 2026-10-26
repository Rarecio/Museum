package config

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseDSN    string `yaml:"database_dsn" env:"DATABASE_DSN" env_required:"true"`
	Port           int    `yaml:"port" env:"PORT" default:"8080"`
	Secret         string `yaml:"secret" env:"SECRET" env_required:"true"`
	ImagesLocation string `yaml:"images_location" env:"IMAGES_LOCATION" env_required:"true"`
}

func Load() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("путь к конфигу не указан")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("конфиг файл не найден: " + configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("ошибка чтения конфига: " + err.Error())
	}
	if cfg.ImagesLocation == "" {
		panic("папка для хранения изображений не указана")
	}
	return &cfg
}

func fetchConfigPath() string {
	var pathFromParam string
	flag.StringVar(&pathFromParam, "config", "", "config file path")
	flag.Parse()
	if pathFromParam != "" {
		return pathFromParam
	}

	pathFromEnv := os.Getenv("CONFIG_PATH")
	if pathFromEnv != "" {
		return pathFromEnv
	}

	return filepath.Join(executableDir(), "config.yml")
}

func executableDir() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	// Разрешаем возможные символьные ссылки на пути к бинарнику
	realExe, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return ""
	}
	return filepath.Dir(realExe)
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
