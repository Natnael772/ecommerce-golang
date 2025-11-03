package uploader

import (
	"fmt"
	"os"
)

func init() {
    if err := os.MkdirAll(UploadPath, os.ModePerm); err != nil {
        panic(fmt.Sprintf("failed to create upload dir: %v", err))
    }
}
