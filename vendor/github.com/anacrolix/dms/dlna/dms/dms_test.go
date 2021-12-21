package dms

import (
	"bytes"
	"net/http"
	"runtime"
	"testing"
)

type safeFilePathTestCase struct {
	root, given, expected string
}

func TestSafeFilePath(t *testing.T) {
	var cases []safeFilePathTestCase
	if runtime.GOOS == "windows" {
		cases = []safeFilePathTestCase{
			{"c:", "/", "c:."},
			{"c:", "/test", "c:test"},
			{"c:\\", "/", "c:\\"},
			{"c:\\", "/test", "c:\\test"},
			{"c:\\hello", "../windows", "c:\\hello\\windows"},
			{"c:\\hello", "/../windows", "c:\\hello\\windows"},
			{"c:\\hello", "/", "c:\\hello"},
			{"c:\\hello", "./world", "c:\\hello\\world"},
			{"c:\\hello", "/", "c:\\hello"},
			// These two ones are invalid but, as this actually prevents to serve them, it is fine
			{"c:\\foo", "c:/windows/", "c:\\foo\\c:\\windows"},
			{"c:\\foo", "e:/", "c:\\foo\\e:"},
		}
	} else {
		cases = []safeFilePathTestCase{
			{"/", "..", "/"},
			{"/hello", "..//", "/hello"},
			{"", "/precious", "precious"},
			{".", "///precious", "precious"},
		}
	}
	t.Logf("running %d test cases", len(cases))
	for _, _case := range cases {
		a := safeFilePath(_case.root, _case.given)
		if a != _case.expected {
			t.Errorf("expected %q from %q and %q but got %q", _case.expected, _case.root, _case.given, a)
		}
	}
}

func TestRequest(t *testing.T) {
	resp, err := http.NewRequest("NOTIFY", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	buf := bytes.NewBuffer(nil)
	resp.Write(buf)
	t.Logf("%q", buf.String())
}

func TestResponse(t *testing.T) {
	var resp http.Response
	resp.StatusCode = http.StatusOK
	resp.Header = make(http.Header)
	resp.Header["SID"] = []string{"uuid:1337"}
	var buf bytes.Buffer
	resp.Write(&buf)
	t.Logf("%q", buf.String())
}
