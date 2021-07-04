package entity

type ConfigRepository interface {
	Create(Config) (Config, error)
	Get() (Config, error)
}

type Config struct {
	Id        int
	Timeout   int64
	LabelName string
}

func NewConfig(timeout int64, labelName string) *Config {
	return &Config{
		Id:        -1,
		Timeout:   timeout,
		LabelName: labelName,
	}
}
