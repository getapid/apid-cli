package command

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/getapid/apid-cli/common/log"
	"github.com/getapid/apid-cli/common/variables"
)

type Executor interface {
	Exec(string, variables.Variables) ([]byte, error)
}

type ShellExecutor struct{}

func NewShellExecutor() Executor {
	return &ShellExecutor{}
}

func (e *ShellExecutor) Exec(command string, vars variables.Variables) ([]byte, error) {
	if len(command) == 0 {
		return []byte{}, errors.New("empty command")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	shell := os.Getenv("SHELL")
	if shell == "" {
		log.L.Warn("SHELL env var not set, using /bin/sh by default")
		shell = "/bin/sh"
	}

	cmd := exec.CommandContext(ctx, shell, "-c", command)
	cmd.Env = append(os.Environ(), getEnvFromVars(vars)...)
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	res := out.String()
	end := len(res) - 1
	if end < 0 {
		end = 0
	}
	return []byte(res[:end]), err
}

func getEnvFromVars(vars variables.Variables) []string {
	var result []string
	for key, value := range vars.Raw() {
		result = append(result, flattenVars(strings.ToUpper(key), value)...)
	}
	return result
}

func flattenVars(namespace string, vars interface{}) []string {
	var result []string
	this, err := json.Marshal(vars)
	if err == nil {
		result = []string{fmt.Sprintf("%s=%s", strings.ToUpper(namespace), this)}
	} else {
		log.L.Debug("could not marshall variables: %s", err)
	}
	switch val := vars.(type) {
	case map[string]interface{}:
		for key, value := range val {
			key = strings.Replace(key, "-", "_", -1)
			result = append(result, flattenVars(strings.ToUpper(namespace+"_"+key), value)...)
		}
		return result
	case []interface{}:
		for index, value := range val {
			result = append(result, fmt.Sprintf("%s=%v", strings.ToUpper(namespace+"_"+strconv.Itoa(index)), value))
		}
		return result
	default:
		result = append(result, fmt.Sprintf("%s=%v", namespace, val))
		return result
	}
}
