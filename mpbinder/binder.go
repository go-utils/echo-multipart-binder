package mpbinder

import (
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/go-utils/echo-multipart-binder/util"
	"github.com/labstack/echo/v4"
	"golang.org/x/xerrors"
)

// NewBindFile - constructor
func NewBindFile(b echo.Binder) echo.Binder {
	return util.BindFunc(
		func(i interface{}, c echo.Context) error {
			if err := b.Bind(i, c); err != nil {
				return xerrors.Errorf("failed to bind method: %w", err)
			}

			ctype := c.Request().Header.Get(echo.HeaderContentType)
			if !(strings.HasPrefix(ctype, echo.MIMEApplicationForm) || strings.HasPrefix(ctype, echo.MIMEMultipartForm)) {
				return nil
			}

			form, err := c.MultipartForm()
			if err != nil {
				return xerrors.Errorf("error in MultipartForm method: %w", err)
			}

			if err = echoBindFile(i, form.File); err != nil {
				return xerrors.Errorf("failed to bind file: %w", err)
			}

			return nil
		},
	)
}

var (
	typeMultipartFileHeader      = reflect.TypeOf((*multipart.FileHeader)(nil))
	typeMultipartSliceFileHeader = reflect.TypeOf(([]*multipart.FileHeader)(nil))
)

func echoBindFile(i interface{}, fileMap map[string][]*multipart.FileHeader) error {
	rv := reflect.Indirect(reflect.ValueOf(i))
	if rv.Kind() != reflect.Struct {
		return xerrors.Errorf("bindFile input not is struct pointer, indirect type is %s", rv.Type().String())
	}

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		fv := rv.Field(i)
		if !fv.CanSet() {
			continue
		}

		ft := rt.Field(i)
		switch ft.Type {
		case typeMultipartFileHeader:
			files := getFiles(fileMap, ft.Name, ft.Tag.Get("form"))
			if len(files) > 0 {
				fv.Set(reflect.ValueOf(files[0]))
			}
		case typeMultipartSliceFileHeader:
			files := getFiles(fileMap, ft.Name, ft.Tag.Get("form"))
			if len(files) > 0 {
				fv.Set(reflect.ValueOf(files))
			}
		}
	}

	return nil
}

func getFiles(files map[string][]*multipart.FileHeader, names ...string) []*multipart.FileHeader {
	for _, name := range names {
		if file, ok := files[name]; ok {
			return file
		}
	}

	return nil
}
