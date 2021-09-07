package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type skopeo struct {}

func (s *skopeo) version() (string, error) {
	return s.execCommand("-v")
}

func (s *skopeo) login(user string, pass string, registry string, tlsVerify bool) (string, error) {
	tlsVerifyOption := fmt.Sprintf("--tls-verify=%t", tlsVerify)
	return s.execCommand("login", "-u", user, "-p", pass, tlsVerifyOption, registry)
}

func (s *skopeo) copy(
	srcImage string,
	descImage string,
	srcTlsVerify bool,
	destTlsVerify bool) (string, error) {
	srcImageURI := fmt.Sprintf("docker://%s", srcImage)
	destImageURI := fmt.Sprintf("docker://%s", descImage)
	srcTlsVerifyOption := fmt.Sprintf("--src-tls-verify=%t", srcTlsVerify)
	destTlsVerifyOption := fmt.Sprintf("--dest-tls-verify=%t", destTlsVerify)
	fmt.Sprintf("srcImageURI: %s", srcImageURI)
	fmt.Sprintf("destImageURI: %s", destImageURI)
	fmt.Sprintf("srcTlsVerifyOption: %s", srcTlsVerifyOption)
	fmt.Sprintf("destTlsVerifyOption: %s", destTlsVerifyOption)
	return s.execCommand("copy", srcImageURI, destImageURI, srcTlsVerifyOption, destTlsVerifyOption)
}

func (s *skopeo) execCommand(arg ...string) (string, error) {
	cmd := exec.Command("skopeo", arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}


