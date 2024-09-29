package config

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const ShellToUse = "bash"

func ShellOut(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func GetCurrentSudoUser() (string, error) {
	err, result, _ := ShellOut("echo $SUDO_USER")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result), nil
}

func GetLocalConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	abs, err := filepath.Abs(path.Join(usr.HomeDir, "."+Company))
	if err != nil {
		return "", err
	}
	return abs, nil
}

func CreateLocalConfigDir() (string, error) {
	local, err := GetLocalConfigPath()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(local); os.IsNotExist(err) {
		os.Mkdir(local, 0755)
	}
	return local, nil
}

func CreateConfigFile(prefix, ext string) (*os.File, error) {
	localPath, err := CreateLocalConfigDir()
	if err != nil {
		return nil, err
	}
	cfgName := fmt.Sprintf("%sconfig-%v.%s", prefix, time.Now().UnixNano(), ext)
	cfgFile := path.Join(localPath, cfgName)
	return os.Create(cfgFile)
}

func SaveToLocalConfig(prefix, ext string, file []byte) (string, error) {
	f, err := CreateConfigFile(prefix, ext)
	if err != nil {
		return "", err
	}
	defer f.Sync()
	defer f.Close()
	_, err = f.Write(file)
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}
