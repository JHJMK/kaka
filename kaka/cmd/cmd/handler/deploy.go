package handler

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"kaka/cmd/cmd/config"
	"kaka/cmd/cmd/easyhost"
	"strconv"
	"strings"
	"text/template"
)

// CreateDir 创建目录
func CreateDir(c config.Config) error {
	fmt.Println("start check dir ...")
	for _, host := range c.ManageHost {
		err := func() error {
			client, err := easyhost.NewSSHClient(host.User, host.Password, host.IP, host.Port)
			if err != nil {
				return err
			}
			defer client.Close()
			cmd := "mkdir -p " + c.InstallPath
			err = easyhost.ExecuteCmd(client, cmd)
			if err != nil {
				return err
			}
			cmd = "mkdir -p " + c.KafkaLogDir
			err = easyhost.ExecuteCmd(client, cmd)
			if err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}

// CopyFile 把kaka文件
func CopyFile(c config.Config) error {
	for _, host := range c.ManageHost {
		fmt.Println("start copy file to " + host.IP)
		// kafka
		binary, err := GetInstallationPackageBytes(config.KafkaFullFileName)
		if err != nil {
			return err
		}
		err = easyhost.TransferFile(binary, "/tmp/"+config.KafkaFullFileName, host)
		if err != nil {
			return err
		}
		// 复制zookeeper
		binary, err = GetInstallationPackageBytes(config.ZookeeperFullFileName)
		if err != nil {
			return err
		}
		err = easyhost.TransferFile(binary, "/tmp/"+config.ZookeeperFullFileName, host)
		if err != nil {
			return err
		}
	}
	return nil
}

func DecompressionAndRename(c config.Config) error {
	// 将压缩包解压到kaka二进制的同级data目录下
	for _, host := range c.ManageHost {
		err := func() error {
			fmt.Println("start decompression and rename host: ", host.IP)
			client, err := easyhost.NewSSHClient(host.User, host.Password, host.IP, host.Port)
			if err != nil {
				return err
			}
			defer client.Close()
			// cd /opt/kaka && tar -zxf /tmp/kafka_2.12-2.6.1.tgz && mv kafka_2.12-2.6.1 kafka
			kafkaCmd := "cd " + c.InstallPath + " && tar -zxf /tmp/" + config.KafkaFullFileName + " && mv " + config.KafkaFileName + " kafka"
			res, err := easyhost.ExecuteCmdWithResponse(client, kafkaCmd)
			if err != nil {
				return errors.Wrap(err, kafkaCmd+" execute error,return message: "+res)
			}
			zkCmd := "cd " + c.InstallPath + " && tar -zxf /tmp/" + config.ZookeeperFullFileName + " && mv " + config.ZookeeperFileName + " zookeeper"
			res, err = easyhost.ExecuteCmdWithResponse(client, zkCmd)
			if err != nil {
				return errors.Wrap(err, zkCmd+" execute error,return message: "+res)
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}

const templatePath = "../config/template/"

func ConfigKafka(c config.Config) error {
	for i, host := range c.ManageHost {
		err := func() error {
			fmt.Println("start config kafka host: ", host.IP)
			// 新建客户端
			sshClient, err := easyhost.NewSSHClient(host.User, host.Password, host.IP, host.Port)
			if err != nil {
				return err
			}
			defer sshClient.Close()
			//1,配置zk
			file, err := config.Template.ReadFile("template/zoo.cfg")
			if err != nil {
				return err
			}
			tpl, err := template.New("zoo.cfg").Parse(string(file))
			if err != nil {
				sshClient.Close()
				return err
			}
			var buf bytes.Buffer
			nodeInfo := make(map[string]string)
			for j, h := range c.ManageHost {
				nodeInfo[string(rune(j+1))] = h.IP // 使用 IP 地址，server.1=IP:2888:3888
			}
			err = tpl.Execute(&buf, map[string]interface{}{
				"NodeInfo": nodeInfo,
			})
			if err != nil {
				return err
			}
			err = easyhost.TransferFile(buf.Bytes(), c.InstallPath+"/zookeeper/conf/zoo.cfg", host)
			if err != nil {
				return err
			}
			err = easyhost.TransferFile([]byte(strconv.Itoa(i)), c.InstallPath+"/zookeeper/myid", host)
			if err != nil {
				return err
			}
			//2.配置启动脚本
			file, err = config.Template.ReadFile("template/kafka-server-start.sh")
			if err != nil {
				return err
			}
			tpl, err = template.New("kafka-server-start.sh").Parse(string(file))
			if err != nil {
				return err
			}
			err = tpl.Execute(&buf, map[string]interface{}{
				"MaxHeapSize": "1G",
				"MinHeapSize": "1G",
			})
			if err != nil {
				return err
			}
			err = easyhost.TransferFile(buf.Bytes(), c.InstallPath+"/kafka/bin/kafka-server-start.sh", host)
			if err != nil {
				return err
			}
			//3.配置kafka/config/server.properties
			file, err = config.Template.ReadFile("template/2.6.1.properties")
			if err != nil {
				return err
			}
			tpl, err = template.New("2.6.1.properties").Funcs(funcMap).Parse(string(file))
			if err != nil {
				return err
			}
			ips := make([]string, len(c.ManageHost))
			for _, h := range c.ManageHost {
				ips[i] = h.IP + ":2181"
			}
			zookeeperConnect := strings.Join(ips, ",")
			params := map[string]interface{}{
				"num.network.threads":            3,
				"num.io.threads":                 8,
				"num.partitions":                 3,
				"default.replication.factor":     3,
				"auto.create.topics.enable":      false,
				"auto.leader.rebalance.enable":   true,
				"delete.topic.enable":            true,
				"message.max.bytes":              1000000,
				"log.retention.ms":               259200000,
				"log.retention.bytes":            -1,
				"replica.fetch.max.bytes":        1048576,
				"log.cleanup.policy":             "delete",
				"min.insync.replicas":            1,
				"unclean.leader.election.enable": false,
			}
			err = tpl.Execute(&buf, map[string]interface{}{
				"ZookeeperConnect": zookeeperConnect,
				"LogDir":           c.KafkaLogDir,
				"BrokerId":         i,
				"NodeIp":           host.IP,
				"Params":           params,
			})
			if err != nil {
				return err
			}
			err = easyhost.TransferFile(buf.Bytes(), c.InstallPath+"/kafka/config/server.properties", host)
			if err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}

	return nil
}
