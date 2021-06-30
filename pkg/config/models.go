package config

type Config struct {
	LogPath     string `json:"log_path"`
	MappingPath string `json:"mapping_path"`
	Server      Server `json:"server"`
	Slack       Slack  `json:"slack"`
	Gitlab      Gitlab `json:"gitlab"`
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Slack struct {
	Token string `json:"token"`
}

type Gitlab struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
}
