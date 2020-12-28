package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func CreateFile(name string) {
	d := []byte("")
	ioutil.WriteFile(name, d, 0644)
}
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
func main() {
	args := os.Args[1:]
	if args[0] == "init" {
		fmt.Println("Chose The Language That The Project Is In")
		fmt.Println("1) Python")
		fmt.Println("2) TypeScript")
		var lang string
		fmt.Scanln(&lang)
		if lang == "1" {
			os.Mkdir("pow_pack", 0755)
			CreateFile(".powp")
		} else if lang == "2" {
			os.Mkdir("pow_pack", 0755)
			CreateFile(".powt")
		} else {
			fmt.Println("invaild")
			os.Exit(1)
		}
	}
	if args[0] == "install" {
		DownloadFile("download.zip", args[1])
		os.Mkdir("pow_pack/"+args[2], 0755)
		Unzip("download.zip", "pow_pack/"+args[2])
	}
}
