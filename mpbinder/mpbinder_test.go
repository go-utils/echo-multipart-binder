package mpbinder

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func createRequest(t *testing.T, filenames ...string) (io.Reader, string) {
	t.Helper()

	buf := bytes.Buffer{}
	mpw := multipart.NewWriter(&buf)
	defer mpw.Close()

	for _, name := range filenames {
		_, err := mpw.CreateFormFile(name, name)

		if err != nil {
			t.Fatal(err)
		}
	}

	return &buf, mpw.FormDataContentType()
}

type Embedded struct {
	E *multipart.FileHeader `form:"e"`
}

type request struct {
	Embedded
	Param string                `param:"p"`
	F     *multipart.FileHeader `form:"f"`
}

func handler(c echo.Context) error {
	var r request
	if err := c.Bind(&r); err != nil {
		return err
	}

	if r.Param != "p" {
		return fmt.Errorf("param differed: %s", r.Param)
	}
	if r.F.Filename != "f" {
		return fmt.Errorf("filename differed for F: %s", r.F.Filename)
	}
	if r.E.Filename != "e" {
		return fmt.Errorf("filename differed for E: %s", r.E.Filename)
	}

	return nil
}

func Test_NewBindFile(t *testing.T) {
	e := echo.New()

	e.Binder = NewBindFile(&echo.DefaultBinder{})

	body, contentType := createRequest(t, "f", "e")

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("p")
	c.SetParamValues("p")

	if err := handler(c); err != nil {
		t.Fatal(err)
	}
}
