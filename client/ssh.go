package client

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

func (sc *SSHClient) Exec(flags []string) (int, string) {
	if len(flags) < 1 {
		return sc.exec("pct")
	}

	switch flags[0] {
	case "pull":
		return sc.pull(flags)
	case "push":
		return sc.push(flags)
	default:
		return sc.exec("pct " + strings.Join(flags, " "))
	}
}

func (sc *SSHClient) exec(cmd string) (int, string) {
	session, err := sc.ssh.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer func() {
		_ = session.Close()
	}()

	out, err := session.CombinedOutput(cmd)
	out = bytes.TrimRight(out, " \n")
	code := 0

	var e *ssh.ExitError
	if errors.As(err, &e) {
		code = e.ExitStatus()
	} else if err != nil {
		code = 1
	}

	return code, string(out)
}

// 下载文件
func (sc *SSHClient) pull(flags []string) (int, string) {
	code, tmp := sc.exec("mktemp -d --suffix=pve-alt-agent")
	defer sc.exec("rm -rf " + tmp)

	// pct pull 100 /remote.file ./local.file
	count := len(flags)
	if count != 4 {
		return 1, "not enough arguments"
	}

	subCmd := "pct " + strings.Join(flags[:count-1], " ")
	localPath := flags[count-1]

	name := filepath.Base(localPath)
	code, out := sc.exec(fmt.Sprintf("%s %s/%s", subCmd, tmp, name))
	if code == 0 {
		err := sc.Download(fmt.Sprintf("%s/%s", tmp, name), localPath)
		if err != nil {
			return 1, err.Error()
		}
	}

	return 0, out
}

// 上传文件
func (sc *SSHClient) push(flags []string) (int, string) {
	_, tmp := sc.exec("mktemp -d --suffix=pve-alt-agent")
	defer sc.exec("rm -rf " + tmp)

	// pct push 100 ./local.file /remote.file
	count := len(flags)
	if count != 4 {
		return 1, "not enough arguments"
	}

	subCmd := "pct " + strings.Join(flags[:count-2], " ")
	localPath := flags[count-2]
	remotePath := flags[count-1]

	name := filepath.Base(localPath)
	err := sc.Upload(localPath, fmt.Sprintf("%s/%s", tmp, name))
	if err != nil {
		return 1, err.Error()
	}

	return sc.exec(fmt.Sprintf("%s %s/%s %s", subCmd, tmp, name, remotePath))
}
