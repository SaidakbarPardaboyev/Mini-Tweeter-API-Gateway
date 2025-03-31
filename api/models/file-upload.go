package models

import "mime/multipart"

type File struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type FileUploadReponse struct {
	Storage          string `json:"storage,omitempty"`
	FileNameDisk     string `json:"file_name_disk,omitempty"`
	FileNameDownload string `json:"file_name_download,omitempty"`
	Link             string `json:"link,omitempty"`
	FileSize         int64  `json:"file_size,omitempty"`
}
