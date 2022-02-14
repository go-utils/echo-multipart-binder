package mjbinder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-utils/echo-multipart-binder/mpbinder"
	"github.com/labstack/echo/v4"
)

func createRequest(t *testing.T, body interface{}, filenames ...string) (io.Reader, string) {
	t.Helper()

	buf := bytes.Buffer{}
	mpw := multipart.NewWriter(&buf)
	defer mpw.Close()

	writer, err := mpw.CreatePart(CreateJSONRequestMIMEHeader())

	if err != nil {
		t.Fatal(err)
	}

	if err := json.NewEncoder(writer).Encode(body); err != nil {
		t.Fatal(err)
	}

	for _, name := range filenames {
		_, err := mpw.CreateFormFile(name, name)

		if err != nil {
			t.Fatal(err)
		}
	}

	return &buf, mpw.FormDataContentType()
}

type request struct {
	A int
	B string
	F *multipart.FileHeader `json:"-" form:"filename.json"`
}

func handler(c echo.Context) error {
	var r request
	if err := c.Bind(&r); err != nil {
		return err
	}

	if r.A != 1 {
		return fmt.Errorf("A is %d", r.A)
	}
	if r.B != "string" {
		return fmt.Errorf("B is %s", r.B)
	}
	fmt.Println(r.F)
	if r.F.Filename != "filename.json" {
		return fmt.Errorf("filename differed: %s", r.F.Filename)
	}

	return nil
}

func Test_NewMultipartJsonBinder(t *testing.T) {
	e := echo.New()

	e.Binder = NewMultipartJSONBinder(
		mpbinder.NewMultipartFileBinder(&echo.DefaultBinder{}),
	)

	body, contentType := createRequest(t, request{
		A: 1,
		B: "string",
	}, "filename.json")

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	if err := handler(c); err != nil {
		t.Fatal(err)
	}
}
