package models

type Config struct {
	Db   Db        `yaml:"db"`
	Mq   Mq        `yaml:"mq"`
	SFID SnowFlake `yaml:"snowflake"`
}
type Mysql struct {
	Addr      string `yaml:"addr"`
	DefaultDB string `yaml:"defaultDB"`
	Passwd    string `yaml:"passwd"`
	Port      int    `yaml:"port"`
	User      string `yaml:"user"`
}
type Db struct {
	Mysql Mysql `yaml:"mysql"`
}
type Kafka struct {
	Addr        string `yaml:"addr"`
	Port        string `yaml:"port"`
	GroupID     string `yaml:"group_id"`
	BlanketPeek int    `yaml:"blanket_peak"`
}
type Mq struct {
	Kafka Kafka `yaml:"kafka"`
}

type SnowFlake struct {
	MachineID int64 `yaml:"machineID"`
}
