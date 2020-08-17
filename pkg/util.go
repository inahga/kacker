package kacker

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var globalConf Configuration

// Configuration sets the behavior of this module.
type Configuration struct {
	VerboseLogging bool
	UseInsecureSSL bool
}

// UseConfiguration sets the behavior of this module.
func UseConfiguration(conf *Configuration) {
	globalConf = *conf
}

func hasCommand(cmd string) bool {
	ret := exec.Command(cmd)
	if err := ret.Run(); err != nil {
		return false
	}
	return true
}

func readReadCloser(readCloser io.ReadCloser, f *os.File, c chan struct{}) {
	scanner := bufio.NewScanner(readCloser)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintln(f, line)
	}
	c <- struct{}{}
}

func runCommands(commands []*exec.Cmd) error {
	for _, command := range commands {
		stdout, err := command.StdoutPipe()
		if err != nil {
			return err
		}
		stderr, err := command.StderrPipe()
		if err != nil {
			return err
		}

		stdoutChan := make(chan struct{})
		stderrChan := make(chan struct{})
		go readReadCloser(stdout, os.Stdout, stdoutChan)
		go readReadCloser(stderr, os.Stderr, stderrChan)
		if err = command.Start(); err != nil {
			return err
		}
		<-stdoutChan
		<-stderrChan
		if err = command.Wait(); err != nil {
			return err
		}
	}
	return nil
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

func resolveToTempFile(dir string, template string, resolve func(io.Writer) error) (string, error) {
	f, err := ioutil.TempFile(dir, template)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if globalConf.VerboseLogging {
		log.Printf("Writing to %s\n", f.Name())
	}
	return f.Name(), resolve(f)
}
