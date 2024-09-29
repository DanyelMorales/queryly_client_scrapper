package cron

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	config2 "news_scrapper/pkg/config"
	"news_scrapper/pkg/style"
	"os"
	"path/filepath"
)

func GetPlaceholder(fallback, rule string) string {
	if rule == "*" {
		return fallback
	}
	id := md5.Sum([]byte(rule))
	return fmt.Sprintf(`%x_%s`, id, fallback)
}

func GetCronExecLine(fallback, timeRule, tool, cfg, logRedirect, options string) *Cron {
	if cfg != "" {
		abs, err := filepath.Abs(cfg)
		if err != nil {
			return nil
		}
		cfg = "--cfg " + abs
	}

	if logRedirect != "" {
		abs, err := filepath.Abs(logRedirect)
		if err != nil {
			return nil
		}
		logRedirect = " >> " + abs + " 2>&1"
	}
	cmd := fmt.Sprintf("%s %s %s %s", tool, cfg, options, logRedirect)
	id := GetPlaceholder(fallback, timeRule)

	return &Cron{
		Id:       id,
		Cmd:      cmd,
		TimeRule: timeRule,
		CronLine: fmt.Sprintf(`%s %s #%s`, timeRule, cmd, id),
	}
}

func CreateTmpFile(handler func(f *os.File) error) error {
	file, err := ioutil.TempFile(os.TempDir(), config2.PREFIX)
	if err != nil {
		return err
	}
	defer file.Close()
	defer os.Remove(file.Name())
	err = handler(file)
	if err != nil {
		return err
	}
	return nil
}

func ModifyCrontab(user string, handler func(f *os.File) error) error {
	return CreateTmpFile(func(file *os.File) error {
		_, err := SaveCrontab(user, file.Name()).Exec()
		if err != nil {
			return err
		}
		style.SuccessfulAction("cloned crontab file")
		err = handler(file)
		if err != nil {
			return err
		}
		_, err = InstallCronFile(user, file.Name()).Exec()
		if err != nil {
			return err
		}
		style.SuccessfulAction("saved crontab configuration")
		return nil
	})
}

func UninstallCron(user string, rule *Cron) error {
	exists, err := IsInit(user)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return ModifyCrontab(user, func(f *os.File) error {
		command := fmt.Sprintf(`sed -i '/%s$/d' %s`, rule.Id, f.Name())
		err, _, stderr := config2.ShellOut(command)
		if err != nil {
			style.FailedActionF("%s - %v", stderr, err)
			return err
		}
		style.SuccessfulAction("removed crontab line")
		return nil
	})
}

func ListCron(user, pattern string) error {
	command := fmt.Sprintf(`%s | awk '/%s/'`, ListCrontab(user).Cmd, pattern)
	err, out, _ := config2.ShellOut(command)
	if err != nil {
		return err
	}
	if len(out) == 0 {
		style.SuccessfulAction("no cron tasks found")
	} else {
		style.SuccessfulActionF("listing installed cron tasks: \n%s", out)
	}
	return nil
}

func InstallCron(user string, cronCfg *Cron) error {
	ModifyCrontab(user, func(f *os.File) error {
		_, err := SaveLine(cronCfg.CronLine, f.Name()).Exec()
		if err != nil {
			return err
		}
		style.SuccessfulAction("installed crontab line")
		return nil
	})
	return nil
}

func InitCron(user string) error {
	return CreateTmpFile(func(f *os.File) error {
		_, err := SaveLine("# --- initialized ---", f.Name()).Exec()
		if err != nil {
			return err
		}
		_, err = InstallCronFile(user, f.Name()).Exec()
		if err != nil {
			return err
		}
		style.SuccessfulAction("initialized crontab")
		return nil
	})
}
