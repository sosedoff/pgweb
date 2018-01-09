package cli

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/sosedoff/pgweb/pkg/command"
)

var (
	latestVersion string
	workDir       = "/tmp"
	distName      = fmt.Sprintf("pgweb_%s_%s", runtime.GOOS, runtime.GOARCH)
	zipName       = distName + ".zip"
)

func printErr(err error) {
	fmt.Println(err)
}

func checkUpdate() {

	if options.DisableCheckUpdate {
		return
	}

	var latestRelease struct{ Name string }

	response, err := http.Get("https://api.github.com/repos/sosedoff/pgweb/releases/latest")
	if err != nil {
		printErr(fmt.Errorf("somethin goes wrong while getting latest release: %s", err))
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&latestRelease)
	if err != nil {
		printErr(fmt.Errorf("something goes wrong while decoding GitHub response: %s", err))
	}

	latestVersion = latestRelease.Name

	if latestVersion != command.VERSION {
		var upgrade string

		fmt.Printf("A new version %s is available. Do you want to upgrade? [y/N] ", latestVersion)
		chars, err := fmt.Scanln(&upgrade)
		if err != nil || (chars > 0 && upgrade == "y") {
			fmt.Println("Downloading and installing new version...")
			installUpdate()
		} else {
			printErr(fmt.Errorf("something goes wrong while updating to new version: %s", err))
			return
		}

	}

}

func installUpdate() {
	var platformDownloadUrl = fmt.Sprintf("https://github.com/sosedoff/pgweb/releases/download/v%s/%s",
		latestVersion, zipName)
	err := downloadFromUrl(platformDownloadUrl)
	if err != nil {
		printErr(fmt.Errorf("error occured during donwload - %s", err))
	}

	err = extractZip()
	if err != nil {
		printErr(fmt.Errorf("error occured during extracting - %s", err))
	}

	// installing new binary
	wd, _ := os.Getwd()

	destPath := filepath.Join(wd, os.Args[0])
	srcPath := filepath.Join(workDir, distName, distName)

	data, err := ioutil.ReadFile(srcPath)
	if err != nil {
		printErr(fmt.Errorf("error occured during replacing - %s", err))
	}
	err = ioutil.WriteFile(destPath, data, 0644)
	if err != nil {
		printErr(fmt.Errorf("error occured during replacing - %s", err))
	}

	_, _, err = syscall.StartProcess(os.Args[0], os.Args, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func extractZip() error {
	ar, err := zip.OpenReader(filepath.Join(workDir, zipName))
	defer ar.Close()

	for _, f := range ar.File {

		rc, err := f.Open()

		if err != nil {
			return err
		}

		dest := filepath.Join(workDir, distName)
		os.MkdirAll(dest, f.FileInfo().Mode())

		fileCopy, err := os.Create(filepath.Join(dest, f.Name))
		if err != nil {
			return err
		}

		_, err = io.Copy(fileCopy, rc)
		if err != nil {
			return err
		}

		rc.Close()
		fileCopy.Close()
	}

	return err
}

func downloadFromUrl(url string) error {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(filepath.Join(workDir, fileName))
	if err != nil {
		return err

	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		return err
	}

	fmt.Println(n, "bytes downloaded.")
	return nil
}
