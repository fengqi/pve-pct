package client

import (
	"io"
	"os"
)

func (sc *SSHClient) Upload(localPath string, remotePath string) (err error) {
	local, err := os.Open(localPath)
	if err != nil {
		return
	}
	defer local.Close()

	remote, err := sc.sftp.Create(remotePath)
	if err != nil {
		return
	}
	defer remote.Close()

	_, err = io.Copy(remote, local)
	return
}

func (sc *SSHClient) Download(remotePath string, localPath string) (err error) {
	local, err := os.Create(localPath)
	if err != nil {
		return
	}
	defer local.Close()

	remote, err := sc.sftp.Open(remotePath)
	if err != nil {
		return
	}
	defer remote.Close()

	if _, err = io.Copy(local, remote); err != nil {
		return
	}

	return local.Sync()
}
