package config

type App struct {
	Env       string `mapstructure:"env" json:"env" yaml:"env"`
	Port      string `mapstructure:"port" json:"port" yaml:"port"`
	AppName   string `mapstructure:"app_name" json:"app_name" yaml:"app_name"`
	AppUrl    string `mapstructure:"app_url" json:"app_url" yaml:"app_url"`
	VideoPath string `mapstructure:"video_path" json:"video_path" yaml:"video_path"`
	SrtPath   string `mapstructure:"srt_path" json:"srt_path" yaml:"srt_path"`
}
