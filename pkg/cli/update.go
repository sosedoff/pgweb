package cli

import (
	"archive/zip"
	"context"
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

	"github.com/google/go-github/github"
	"github.com/sosedoff/pgweb/pkg/command"
)

var (
	latestVersion string
	workDir       = "/tmp"
	distName      = fmt.Sprintf("pgweb_%s_%s", runtime.GOOS, runtime.GOARCH)
	zipName       = distName + ".zip"
)

func checkUpdate() {

	if options.CheckUpdate == false {
		return
	}

	client := github.NewClient(nil)

	latestRelease, _, err := client.Repositories.GetLatestRelease(context.Background(), "sosedoff", "pgweb")
	if err != nil {
		fmt.Errorf("GitHub is unavailable")
	}

	latestVersion = *latestRelease.TagName

	if latestVersion != command.VERSION {
		var upgrade string

		fmt.Printf("A new version %s is available. Do you want to upgrade? [y/N]", latestVersion)
		chars, err := fmt.Scanln(&upgrade)
		if err != nil || (chars > 0 && upgrade == "y") {
			fmt.Println("Downloading and installing new version...")
			installUpdate()
		} else {
			return
		}

	}

}

func installUpdate() {
	var platformDownloadUrl = fmt.Sprintf("https://github.com/sosedoff/pgweb/releases/download/%s/%s",
		latestVersion, zipName)
	err := downloadFromUrl(platformDownloadUrl)
	if err != nil {
		fmt.Errorf("error occured during donwload - %s", err)
	}

	err = extractZip()
	if err != nil {
		fmt.Errorf("error occured during extracting - %s", err)
	}

	// installing new binary
	wd, _ := os.Getwd()

	destPath := filepath.Join(wd, os.Args[0])
	srcPath := filepath.Join(workDir, distName, distName)

	data, err := ioutil.ReadFile(srcPath)
	if err != nil {
		fmt.Errorf("error occured during replacing - %s", err)
	}
	err = ioutil.WriteFile(destPath, data, 0644)
	if err != nil {
		fmt.Errorf("error occured during replacing - %s", err)
	}

	_, err = syscall.ForkExec(os.Args[0], os.Args, nil)
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
