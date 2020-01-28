package models


type SystemConfiguration struct {

	Mysql mysql `json:"mysql" yaml:"mysql"`

	Server server `json:"server" yaml:"server"`

	Redis redis `json:"redis" yaml:"redis"`

	Dijan dijan `json:"dijan" yaml:"dijan"`
}


type server struct {

	Salt string `json:"salt" yaml:"salt"`

}

type mysql struct {

	DB string `json:"db" yaml:"db"`

	Host string `json:"host" yaml:"host"`

	Port string `json:"port" yaml:"port"`

	Username string `json:"username" yaml:"username"`

	Password string `json:"password" yaml:"password"`

}

type redis struct {

	StartNodes []string `json:"start_nodes" yaml:"start_nodes"`

	Tcp string `json:"tcp" yaml:"tcp"`

}

type dijan struct {

	Host string `json:"host" yaml:"host"`

	Port string `json:"port" yaml:"port"`
}