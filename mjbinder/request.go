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

// CreateJsonRequestMIMEHeader creates a MIME header for JSON part in multipart requests
func CreateJsonRequestMIMEHeader() textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(JsonPartKey), escapeQuotes(JsonPartKey)))
	h.Set("Content-Type", "application/json")
	h.Set(JsonPartHeaderKey, "1")

	return h
}
