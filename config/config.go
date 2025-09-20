package config

import (
	"embed"
	"os"
	"path/filepath"
)

// Kafka 版本常量
const (
	KafkaFullFileName     = "kafka_2.12-2.6.1.tgz"
	KafkaFileName         = "kafka_2.12-2.6.1"
	ZookeeperFullFileName = "apache-zookeeper-3.5.8-bin.tar.gz"
	ZookeeperFileName     = "apache-zookeeper-3.5.8-bin"
	JdkFullFileName       = "openjdk-8u44-linux-x64.tar.gz"
	JdkFileName           = "java-se-8u44-ri"
)

var StaticFile = map[string]string{
	KafkaFullFileName:     KafkaFileName,
	ZookeeperFullFileName: ZookeeperFileName,
}

//go:embed template/*
var Template embed.FS

type Host struct {
	IP       string `json:"IP"`
	Port     string `json:"Port"`
	User     string `json:"User"`
	Password string `json:"Password"`
}

type Config struct {
	ManageHost  []Host `json:"ManageHost"`
	InstallPath string `json:"InstallPath"`
	DataDir     string `json:"DataDir"`
}

// GetExecuteDir 获取当前执行目录,结尾不带 /
func GetExecuteDir() (string, error) {
	executable, err := os.Executable()
	return filepath.Dir(executable), err
}
