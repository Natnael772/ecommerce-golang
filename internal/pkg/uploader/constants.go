package uploader

const (
    MaxUploadSize = 10 << 20 // 10MB
    UploadPath    = "./uploads"
)

type ctxKey string

const (
    UploadedFileKey  ctxKey = "uploadedFile"
    UploadedFilesKey ctxKey = "uploadedFiles"
)
