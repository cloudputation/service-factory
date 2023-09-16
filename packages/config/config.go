package config


type Config struct {
    Nomad   Nomad  `hcl:"nomad,block"`
    Server  Server `hcl:"repo,block"`
}

type Server struct {
    ServerPort    string `hcl:"port"`
    ServerAddress string `hcl:"address"`
}

type Nomad struct {
    NomadAddress string `hcl:"server_address"`
}


var terraformDir = "terraform/"
