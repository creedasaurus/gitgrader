package main

import (
	"os"
	"log"
	//"fmt"
	"path/filepath"
	"regexp"
	"archive/zip"
	"io"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough args")
	}
	targetDir := os.Args[1]
	endingDir := os.Args[2]
	ignoreFile := os.Args[3]

	submissionFileRegex := regexp.MustCompile(`^([a-z]*)_?(late)?_([0-9]+)_([0-9]+)_(.*)(\..*)$`)

	err := os.Mkdir(endingDir, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		ms := submissionFileRegex.FindStringSubmatch(info.Name())
		if info.Name() == ignoreFile {
			log.Println(info.Name())
		}

		if len(ms) < 6 {
			log.Println("Not a match")
			return nil
		}

		os.Mkdir(endingDir+"/"+ms[1], os.ModePerm)

		if ms[6] == ".zip" {
			r, err := zip.OpenReader(path)
			if err != nil {
				log.Println(err)
			}

			for _, f := range r.File {

				if f.FileInfo().Name() == ignoreFile {
					log.Println(f.FileInfo().Name())
					continue
				}

				of, err := f.Open()
				if err != nil {
					return err
				}
				defer of.Close()

				nf, err := os.OpenFile(endingDir+"/"+ms[1]+"/"+f.FileInfo().Name(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				if err != nil {
					return err
				}
				defer nf.Close()

				_, err = io.Copy(nf, of)
				if err != nil {
					return err
				}

			}

		} else {
			os.Rename(targetDir+"/"+info.Name(), endingDir+"/"+ms[1]+"/"+info.Name())
		}

		return nil
	})
}
