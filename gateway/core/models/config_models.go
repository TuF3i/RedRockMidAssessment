package models

type Config struct {
	Db        Db        `yaml:"db"`
	Mq        Mq        `yaml:"mq"`
	HertzAPI  HertzAPI  `yaml:"hertzApi"`
	SnowFlake SnowFlake `json:"snowflake"`
}

type SnowFlake struct {
	MachineID int64 `yaml:"machineID"`
}

type Mysql struct {
	Addr      string `yaml:"addr"`
	DefaultDB string `yaml:"defaultDB"`
	Passwd    string `yaml:"passwd"`
	Port      int    `yaml:"port"`
	User      string `yaml:"user"`
}
type Redis struct {
	Addr      string `yaml:"addr"`
	DefaultDB int    `yaml:"defaultDB"`
	Passwd    string `yaml:"passwd"`
	Port      int    `yaml:"port"`
}
type Db struct {
	Mysql Mysql `yaml:"mysql"`
	Redis Redis `yaml:"redis"`
}
type HertzAPI struct {
	ListenAddr  string `yaml:"listenAddr"`
	ListenPort  int    `yaml:"listenPort"`
	MonitorPort int    `yaml:"monitorPort"`
}

type Mq struct {
	Kafka Kafka `yaml:"kafka"`
}

type Kafka struct {
	Addr     string `yaml:"addr"`
	Port     string `yaml:"port"`
	ClientID string `yaml:"client_id"`
}
