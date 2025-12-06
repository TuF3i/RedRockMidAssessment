package models

type Config struct {
	Db        Db        `yaml:"db"`
	Mq        Mq        `yaml:"mq"`
	Snowflake Snowflake `yaml:"snowflake"`
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
type Kafka struct {
	Addr        string `yaml:"addr"`
	BlanketPeak int    `yaml:"blanket_peak"`
	Port        string `yaml:"port"`
}
type Mq struct {
	Kafka Kafka `yaml:"kafka"`
}
type Snowflake struct {
	MachineID int64 `yaml:"machineID"`
}
