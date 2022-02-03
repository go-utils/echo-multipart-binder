package mjbinder

import (
	"fmt"
	"net/textproto"
	"strings"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// CreateJSONRequestMIMEHeader creates a MIME header for JSON part in multipart requests
func CreateJSONRequestMIMEHeader() textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(JSONPartKey), escapeQuotes(JSONPartKey)))
	h.Set("Content-Type", "application/json")

	return h
}
