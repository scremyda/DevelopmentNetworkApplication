package utils

import (
	"fmt"
	"github.com/rs/xid"
	"strings"
)

func GenerateUniqueName(imageName *string) error {
	parts := strings.Split(*imageName, ".")
	if len(parts) > 1 {
		fileExt := parts[len(parts)-1]
		uniqueID := xid.New()
		*imageName = fmt.Sprintf("%s.%s", uniqueID.String(), fileExt)
		return nil
	}
	return fmt.Errorf("uncorrect file name. not fount image extension")
}
