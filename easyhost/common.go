package easyhost

import (
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

func NewSSHClient(user, password, ip, port string) (*ssh.Client, error) {
	// 创建ssh客户端
	sshconfig := ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	return ssh.Dial("tcp", net.JoinHostPort(ip, port), &sshconfig)
}

func ExecuteCmdWithResponse(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	output, err := session.CombinedOutput(cmd)
	return string(output), err
}

func ExecuteCmd(client *ssh.Client, cmd string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	return session.Run(cmd)
}
