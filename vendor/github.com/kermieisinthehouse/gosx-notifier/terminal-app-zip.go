package gosxnotifier

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	zipPath        = "terminal-notifier.temp.zip"
	executablePath = "terminal-notifier.app/Contents/MacOS/terminal-notifier"
	tempDirSuffix  = "gosxnotifier"
)

var (
	rootPath  string
	FinalPath string
)

func supportedOS() bool {
	if runtime.GOOS == "darwin" {
		return true
	} else {
		log.Print("OS does not support terminal-notifier")
		return false
	}
}

func init() {
	if supportedOS() {
		err := installTerminalNotifier()
		if err != nil {
			log.Fatalf("Could not install Terminal Notifier to a temp directory: %s", err)
		} else {
			FinalPath = filepath.Join(rootPath, executablePath)
		}
	}
}

func exists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func installTerminalNotifier() error {
	rootPath = filepath.Join(os.TempDir(), tempDirSuffix)

	//if terminal-notifier.app already installed no-need to re-install
	if exists(filepath.Join(rootPath, executablePath)) {
		return nil
	}
	buf := bytes.NewReader(terminalnotifier())
	reader, err := zip.NewReader(buf, int64(buf.Len()))
	if err != nil {
		return err
	}
	err = unpackZip(reader, rootPath)
	if err != nil {
		return fmt.Errorf("could not unpack zip terminal-notifier file: %s", err)
	}

	err = os.Chmod(filepath.Join(rootPath, executablePath), 0755)
	if err != nil {
		return fmt.Errorf("could not make terminal-notifier executable: %s", err)
	}

	return nil
}

func unpackZip(reader *zip.Reader, tempPath string) error {
	for _, zipFile := range reader.File {
		name := zipFile.Name
		mode := zipFile.Mode()
		if mode.IsDir() {
			if err := os.MkdirAll(filepath.Join(tempPath, name), 0755); err != nil {
				return err
			}
		} else {
			if err := unpackZippedFile(name, tempPath, zipFile); err != nil {
				return err
			}
		}
	}

	return nil
}

func unpackZippedFile(filename, tempPath string, zipFile *zip.File) error {
	writer, err := os.Create(filepath.Join(tempPath, filename))

	if err != nil {
		return err
	}

	defer writer.Close()

	reader, err := zipFile.Open()
	if err != nil {
		return err
	}

	defer reader.Close()

	if _, err = io.Copy(writer, reader); err != nil {
		return err
	}

	return nil
}
