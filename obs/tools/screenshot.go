package tools

import (
	"encoding/base64"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func screenshotImage(imageData, imageFormat string) mcp.Content {
	mime := "image/png"
	if imageFormat != "" {
		mime = "image/" + strings.ToLower(imageFormat)
	}
	raw := imageData
	if rest, ok := strings.CutPrefix(imageData, "data:"); ok {
		if med, data, found := strings.Cut(rest, ","); found {
			raw = data
			if mt, _, ok := strings.Cut(med, ";"); ok && mt != "" {
				mime = mt
			}
		}
	}
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		decoded = []byte(raw)
	}
	return &mcp.ImageContent{Data: decoded, MIMEType: mime}
}
