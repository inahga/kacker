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
	"os/signal"
	"syscall"
)

var globalConf Configuration

// Configuration sets the behavior of this module.
type Configuration struct {
	VerboseLogging bool
	NoVerifySSL    bool
	KeepTempFiles  bool
	RelativeDir    string
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
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	var target *exec.Cmd
	go func() {
		<-ch
		target.Process.Signal(syscall.SIGTERM)
	}()

	for _, command := range commands {
		target = command

		stdout, err := target.StdoutPipe()
		if err != nil {
			return err
		}
		stderr, err := target.StderrPipe()
		if err != nil {
			return err
		}
		stdoutChan := make(chan struct{})
		stderrChan := make(chan struct{})
		go readReadCloser(stdout, os.Stdout, stdoutChan)
		go readReadCloser(stderr, os.Stderr, stderrChan)

		if err = target.Start(); err != nil {
			return err
		}
		<-stdoutChan
		<-stderrChan
		if err = target.Wait(); err != nil {
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
	if globalConf.NoVerifySSL {
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
	ret := f.Name()
	if err != nil {
		return "", err
	}

	func() {
		defer func() {
			f.Close()
			if r := recover(); r != nil && !globalConf.KeepTempFiles {
				log.Println(r)
				rmErr := os.Remove(ret)
				if rmErr != nil {
					log.Printf("Recovering: Could not delete temp file %s\n", ret)
				}
			}
		}()
		if globalConf.VerboseLogging {
			log.Printf("Writing to %s\n", f.Name())
		}
		if err = resolve(f); err != nil {
			panic(err)
		}
	}()
	return ret, err
}

func stringSliceContains(arr []string, s string) bool {
	for _, elem := range arr {
		if elem == s {
			return true
		}
	}
	return false
}
