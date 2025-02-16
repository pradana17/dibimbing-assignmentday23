package models

type CreateDirectoryRequest struct {
	DirectoryName string `json:"directory_name" binding:"required"`
}

type CreateFileRequest struct {
	FileName      string `json:"file_name" binding:"required"`
	DirectoryName string `json:"directory_name" binding:"required"`
	Content       string `json:"content" binding:"required"`
}

type ReadFileRequest struct {
	FileName      string `json:"file_name" binding:"required"`
	DirectoryName string `json:"directory_name" binding:"required"`
}

type RenameFileRequest struct {
	FileName      string `json:"file_name" binding:"required"`
	NewFileName   string `json:"new_file_name" binding:"required"`
	DirectoryName string `json:"directory_name" binding:"required"`
}

type FileUploadResponse struct {
	FileName string `json:"file_name" binding:"required"`
	Message  string `json:"message" binding:"required"`
	Path     string `json:"path" binding:"required"`
}

type DownloadFileRequest struct {
	FileName      string `json:"file_name" binding:"required"`
	DirectoryName string `json:"directory_name" binding:"required"`
}
