package mjbinder

import (
	"encoding/json"
	"mime/multipart"
	"strings"

	"github.com/go-utils/echo-multipart-binder/util"
	"github.com/labstack/echo/v4"
	"golang.org/x/xerrors"
)

// JsonPartKey is the key for a JSON part in the multipart request
const JsonPartKey = "x-multipart-json-binder-request-json"

// JsonPartHeaderKey is the key in headers for a JSON part that must be set to "1"
const JsonPartHeaderKey = "X-Multipart-Json-Request"

// NewMultipartJsonBinder can bind JSON fields in multipart data
func NewMultipartJsonBinder(b echo.Binder) echo.Binder {
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

			files, ok := form.File[JsonPartKey]

			if !ok || len(files) == 0 {
				return nil
			}

			if err = bindJsonPart(i, files[0]); err != nil {
				return xerrors.Errorf("failed to bind file: %w", err)
			}

			return nil
		},
	)
}

func bindJsonPart(i interface{}, file *multipart.FileHeader) error {
	if file.Header.Get("Content-Type") != "application/json" {
		return nil
	}

	if file.Header.Get("X-Multipart-Json-Request") != "1" {
		return nil
	}

	fp, err := file.Open()

	if err != nil {
		return xerrors.Errorf("failed to open the multipart stream: %w", err)
	}

	defer fp.Close()

	if err := json.NewDecoder(fp).Decode(i); err != nil {
		return xerrors.Errorf("failed to parse a file as JSON: %w", err)
	}

	return nil
}
