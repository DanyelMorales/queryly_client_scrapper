package cron

import (
	"errors"
	"fmt"
	"news_scrapper/pkg/config"
	"strings"
)

type Wrapper struct {
	Cmd string
}

type Cron struct {
	Id       string
	Cmd      string
	TimeRule string
	CronLine string
}

func GetCrontabCmd(user string) Wrapper {
	if len(user) > 0 {
		user = "-u " + user
	}
	cmd := fmt.Sprintf("crontab %s ", user)
	return Wrapper{Cmd: cmd}
}

func ListCrontab(user string) Wrapper {
	cmd := GetCrontabCmd(user)
	cmdStr := fmt.Sprintf("%s -l", cmd.Cmd)
	return Wrapper{Cmd: cmdStr}
}

func SaveLine(line, cronFile string) Wrapper {
	return Wrapper{Cmd: fmt.Sprintf(`echo "%s" >> %s`, line, cronFile)}
}

func SaveCrontab(user, file string) Wrapper {
	cmd := ListCrontab(user)
	cmdStr := fmt.Sprintf("%s > %s", cmd.Cmd, file)
	return Wrapper{Cmd: cmdStr}
}

func InstallCronFile(user, file string) Wrapper {
	cmd := GetCrontabCmd(user).Cmd
	cmd = fmt.Sprintf("%s %s", cmd, file)
	return Wrapper{Cmd: cmd}
}

func IsInit(user string) (bool, error) {
	cmd := ListCrontab(user)
	_, err := cmd.Exec()
	if err != nil {
		if strings.Contains(err.Error(), "no crontab for") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c Wrapper) Exec() (string, error) {
	err, std, stderr := config.ShellOut(c.Cmd)
	if len(stderr) > 0 {
		return "", errors.New(stderr)
	}
	if err != nil {
		return "", err
	}

	return std, nil
}
