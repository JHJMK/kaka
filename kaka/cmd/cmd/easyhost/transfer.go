package easyhost

import (
	"fmt"
	"github.com/pkg/sftp"
	"kaka/cmd/cmd/config"
)

// TransferFile 文件传输或写入文件
func TransferFile(copyFileData []byte, targetFilePath string, host config.Host) error {
	// 创建ssh客户端
	client, err := NewSSHClient(host.User, host.Password, host.IP, host.Port)
	if err != nil {
		return err
	}
	// 创建sftp客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	// 创建远程文件
	remoteFile, err := sftpClient.Create(targetFilePath)
	if err != nil {
		return fmt.Errorf("can not create remote file: %v", err)
	}
	defer remoteFile.Close()
	// 写入数据
	_, err = remoteFile.Write(copyFileData)
	if err != nil {
		return fmt.Errorf("can not write remote file: %v", err)
	}
	return nil
}
