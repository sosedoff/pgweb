package service

import (
	"bytes"
	"errors"
	"os/exec"
	"text/template"
)

func buildURL(url string, address addressInfo) (string, error) {
	urlTemplate, err := template.New("URL").Parse(url)
	if err != nil {
		return "", err
	}
	urlBuffer := &bytes.Buffer{}
	if err := urlTemplate.Execute(urlBuffer, address); err != nil {
		return "", err
	}
	return urlBuffer.String(), nil
}

func buildCommand(arguments []string, address addressInfo) (*exec.Cmd, error) {
	if len(arguments) == 0 {
		return nil, errors.New("empty command")
	}

	command := []string{}
	for _, argument := range arguments {
		argTemplate, err := template.New("command").Parse(argument)
		if err != nil {
			return nil, err
		}

		argBuffer := &bytes.Buffer{}
		if err := argTemplate.Execute(argBuffer, address); err != nil {
			return nil, err
		}
		command = append(command, argBuffer.String())
	}

	return exec.Command(command[0], command[1:]...), nil
}
