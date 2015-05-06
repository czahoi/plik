/**

    Plik upload client

The MIT License (MIT)

Copyright (c) <2015>
	- Mathieu Bodjikian <mathieu@bodjikian.fr>
	- Charles-Antoine Mathieu <skatkatt@root.gg>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
**/

package zip

import (
	"io"
	"os/exec"
	//	"strings"
	"errors"
	"fmt"
	"github.com/root-gg/utils"
	"os"
	"path/filepath"
	"strings"
)

type ZipBackendConfig struct {
	Zip     string
	Options string
}

func NewZipBackendConfig(config map[string]interface{}) (this *ZipBackendConfig) {
	this = new(ZipBackendConfig)
	this.Zip = "/bin/zip"
	utils.Assign(this, config)
	return
}

type ZipBackend struct {
	Config *ZipBackendConfig
}

func NewZipBackend(config map[string]interface{}) (this *ZipBackend, err error) {
	this = new(ZipBackend)
	this.Config = NewZipBackendConfig(config)
	if _, err := os.Stat(this.Config.Zip); os.IsNotExist(err) || os.IsPermission(err) {
		if this.Config.Zip, err = exec.LookPath("zip"); err != nil {
			err = errors.New("zip binary not found in $PATH, please install or edit ~/.plickrc")
		}
	}
	return
}

func (this *ZipBackend) Configure(arguments map[string]interface{}) (err error) {
	if arguments["--archive-options"] != nil && arguments["--archive-options"].(string) != "" {
		this.Config.Options = arguments["--archive-options"].(string)
	}
	return
}

func (this *ZipBackend) Archive(files []string, writer io.WriteCloser) (name string, err error) {
	if len(files) == 0 {
		fmt.Println("Unable to make a zip archive from STDIN")
		os.Exit(1)
		return
	}

	name = "archive"
	if len(files) == 1 {
		name = filepath.Base(files[0])
	}
	name += ".zip"

	args := make([]string, 0)
	args = append(args, strings.Fields(this.Config.Options)...)
	args = append(args, "-r", "-")
	args = append(args, files...)

	cmd := exec.Command(this.Config.Zip, args...)
	cmd.Stdout = writer
	cmd.Stderr = os.Stderr
	go func() {
		err = cmd.Start()
		if err != nil {
			fmt.Printf("Unable to run zip cmd : %s\n", err)
			os.Exit(1)
			return
		}
		err = cmd.Wait()
		if err != nil {
			fmt.Printf("Unable to run zip cmd : %s\n", err)
			os.Exit(1)
			return
		}
		err = writer.Close()
		if err != nil {
			fmt.Printf("Unable to run zip cmd : %s\n", err)
			return
		}
	}()
	return
}

func (this *ZipBackend) Comments() string {
	return ""
}

func (this *ZipBackend) GetConfiguration() interface{} {
	return this.Config
}