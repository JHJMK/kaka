package handler

import (
	"fmt"
	"kaka/config"
	"kaka/easyhost"
)

func Clear(c config.Config) error {
	for _, host := range c.ManageHost {
		fmt.Println("clear " + host.IP)
		err := func() error {
			sshClient, err := easyhost.NewSSHClient(host.User, host.Password, host.IP, host.Port)
			if err != nil {
				return err
			}
			defer sshClient.Close()
			cmd := "rm -rf " + c.InstallPath
			err = easyhost.ExecuteCmd(sshClient, cmd)
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
