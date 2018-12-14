package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os/exec"
)

func RunCommand(cmds ...string) error {
	cmd := exec.Command(cmds[0], cmds[1:]...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func GenerateMd5CipherString(originalText string) string {
	h := md5.New()
	h.Write([]byte(originalText))
	return hex.EncodeToString(h.Sum(nil))
}
