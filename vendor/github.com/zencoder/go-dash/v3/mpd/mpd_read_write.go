package mpd

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"
	"os"
)

// Reads an MPD XML file from disk into a MPD object.
// path - File path to an MPD on disk
func ReadFromFile(path string) (*MPD, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Read(f)
}

// Reads a string into a MPD object.
// xmlStr - MPD manifest data as a string.
func ReadFromString(xmlStr string) (*MPD, error) {
	b := bytes.NewBufferString(xmlStr)
	return Read(b)
}

// Reads from an io.Reader interface into an MPD object.
// r - Must implement the io.Reader interface.
func Read(r io.Reader) (*MPD, error) {
	var mpd MPD
	d := xml.NewDecoder(r)
	err := d.Decode(&mpd)
	if err != nil {
		return nil, err
	}
	return &mpd, nil
}

// Writes an MPD object to a file on disk.
// path - Output path to write the manifest to.
func (m *MPD) WriteToFile(path string) error {
	// Open the file to write the XML to
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = m.Write(f); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	return err
}

// Writes an MPD object to a string.
func (m *MPD) WriteToString() (string, error) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err := m.Write(w)
	if err != nil {
		return "", err
	}
	err = w.Flush()
	if err != nil {
		return "", err
	}
	return b.String(), err
}

// Writes an MPD object to an io.Writer interface
// w - Must implement the io.Writer interface.
func (m *MPD) Write(w io.Writer) error {
	b, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	_, _ = w.Write([]byte(xml.Header))
	_, _ = w.Write(b)
	_, _ = w.Write([]byte("\n"))
	return nil
}
