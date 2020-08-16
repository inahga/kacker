package kacker

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

func hasCommand(cmd string) bool {
	ret := exec.Command(cmd)
	if err := ret.Run(); err != nil {
		return false
	}
	return true
}

// HasPacker checks whether packer is installed on the system.
func HasPacker() bool {
	return hasCommand("packer")
}

// HasKsvalidator checks whether ksvalidator, as part of pykickstart, is
// installed on the system.
func HasKsvalidator() bool {
	return hasCommand("ksvalidator")
}

func escape(str string) string {
	str = fmt.Sprintf("%q", str)
	return str[1 : len(str)-1]
}

func resolveFragment(fragment string) (string, error) {
	ret, err := ioutil.ReadFile(fragment)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func resolveURL(url string) (string, error) {
	var resp *http.Response
	var err error

	if globalConf.UseInsecureSSL {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err = client.Get(url)
	} else {
		resp, err = http.Get(url)
	}

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}
