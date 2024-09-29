package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"news_scrapper/pkg/style"
	"os"
	"path/filepath"
)

type PortalInterface int

const (
	Queryly PortalInterface = iota
	CustomProvider
)

type Scrapper struct {
	OutputPath string `json:"outputPath"`
}
type PortalOption struct {
	Context, Host, ApiKey, Selector string
	Enabled, OverrideExistingNews   bool
	Type                            PortalInterface
}

type RobotoMode string

const (
	PREFIX = "newsbot-"
)

const (
	MasterMode = "master"
	SlaveMode  = "slave"
)

type Roboto struct {
	Mode           RobotoMode     `json:"mode"`
	PortalRegistry []PortalOption `json:"registry"`
	ScrapperConfig Scrapper       `json:"scrapperConfig"`
}

type Settings struct {
	Roboto Roboto `json:"roboto"`
}

func LoadContextSettings(fileName string, s *Settings) error {
	contextPath, _ := os.Getwd()
	return LoadSettingsFile(contextPath+"/"+fileName, s)
}

func LoadSettingsFile(path string, s *Settings) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return ToSettingsFile(b, s)
}

func ToSettingsFile(b []byte, s *Settings) error {
	err := json.Unmarshal(b, s)
	if err != nil {
		style.FailedAction("configuration schema is not valid")
		return err
	}
	return nil
}

func SaveSettingsFromStdin() (string, error) {
	output, err := getStdin()
	if err != nil {
		return "", err
	}
	if len(output) == 0 {
		return "", errors.New("pipe cannot process an empty input")
	}
	f, err := CreateConfigFile(PREFIX, "json")
	if err != nil {
		return "", err
	}
	defer f.Sync()
	defer f.Close()

	f.Write(output)
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func AutoResolveSettings(stdin bool) Settings {
	settings := Settings{}
	var err error
	if stdin {
		err = LoadSettingsFromStdin(&settings)
	} else {
		abs, errAbs := filepath.Abs(viper.GetString(ConfigurationFile))
		err = errAbs
		if err == nil {
			err = LoadContextSettings(abs, &settings)
		}
	}
	if err != nil {
		style.FailedActionF("%v", err)
		os.Exit(1)
	}
	return settings
}

func LoadSettingsFromStdin(s *Settings) error {
	output, err := getStdin()
	if err != nil {
		return err
	}
	return ToSettingsFile(output, s)
}

func (m *RobotoMode) UnmarshalJSON(b []byte) error {
	type LT RobotoMode
	var r = (*LT)(m)
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}
	switch *m {
	case MasterMode, SlaveMode:
		return nil
	}
	return errors.New("invalid roboto mode")
}

func getStdin() ([]byte, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if info.Mode()&os.ModeCharDevice != 0 {
		return nil, errors.New("there is no data on stdin, flag " + CfgStdinFlag + " is intended to work with pipes")
	}

	reader := bufio.NewReader(os.Stdin)
	var output []byte
	for {
		input, err := reader.ReadByte()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	return output, nil
}
