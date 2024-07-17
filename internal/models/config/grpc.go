package config

type GrpcConfig struct {
	Host     string `json:"host"`
	HttpPort string `json:"httpPort"`
	GrpcPort string `json:"grpcPort"`
}