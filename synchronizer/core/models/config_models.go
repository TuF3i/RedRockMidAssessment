package models

type Config struct {
	Db        Db        `yaml:"db"`
	Mq        Mq        `yaml:"mq"`
	Snowflake Snowflake `yaml:"snowflake"`
}

type Redis struct {
	Addr      string `yaml:"addr"`
	DefaultDB int    `yaml:"defaultDB"`
	Passwd    string `yaml:"passwd"`
	Port      int    `yaml:"port"`
}
type Db struct {
	Redis Redis `yaml:"redis"`
}
type Kafka struct {
	Addr        string `yaml:"addr"`
	BlanketPeak int    `yaml:"blanketPeak"`
	Port        string `yaml:"port"`
}
type Mq struct {
	Kafka Kafka `yaml:"kafka"`
}
type Snowflake struct {
	MachineID int64 `yaml:"machineID"`
}
