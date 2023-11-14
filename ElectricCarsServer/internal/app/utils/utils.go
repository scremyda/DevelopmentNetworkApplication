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

const (
	Query = `
        SELECT assemblies.*, autoparts.*
        FROM assemblies
        JOIN autopart_assemblies ON autopart_assemblies.assembly_id = assemblies.id
        JOIN autoparts ON autoparts.id = autopart_assemblies.autopart_id
        WHERE autopart_assemblies.assembly_id = ? AND autopart_assemblies.count = 1
        AND assemblies.deleted_at IS NULL
        LIMIT 1
    `
)
