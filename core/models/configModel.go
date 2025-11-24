package models

type Config struct {
	Db       Db       `yaml:"db"`
	HertzAPI HertzAPI `yaml:"hertzApi"`
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
	ListenAddr string `yaml:"listenAddr"`
	ListenPort int    `yaml:"listenPort"`
}
