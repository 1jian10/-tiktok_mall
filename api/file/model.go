package file

import "mall/api"

type UploadResp struct {
	Status   api.Status `json:"status"`
	FilePath string     `json:"file_path"`
}

type Name struct {
	FileName string `uri:"filename" binding:"required"`
}
