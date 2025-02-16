package controllers

import (
	"assignmentday23/models"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileController struct {
	DB            *gorm.DB
	downloadMutex sync.Mutex
}

func NewFileController(db *gorm.DB) *FileController {
	return &FileController{
		DB: db,
	}
}

func (fc *FileController) CreateDirectory(c *gin.Context) {
	var req models.CreateDirectoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := os.Mkdir(req.DirectoryName, 0755)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Directory created successfully"})
}

func (fc *FileController) CreateFile(c *gin.Context) {
	var req models.CreateFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := os.MkdirAll(req.DirectoryName, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filepath := filepath.Join(req.DirectoryName, req.FileName)

	file, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer file.Close()

	_, err = file.WriteString(req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File created successfully", "path": filepath})
}

func (fc *FileController) ReadFile(c *gin.Context) {

	var req models.ReadFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filepath := filepath.Join(req.DirectoryName, req.FileName)
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": string(data)})
}

func (fc *FileController) RenameFile(c *gin.Context) {
	var req models.RenameFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	oldPath := filepath.Join(req.DirectoryName, req.FileName)
	newPath := filepath.Join(req.DirectoryName, req.NewFileName)

	err := os.Rename(oldPath, newPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File renamed successfully"})
}

func (fc *FileController) UploadFile(c *gin.Context) {
	const uploadDir = "uploads"

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := filepath.Base(file.Filename)
	dst := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.FileUploadResponse{
		Message:  "File uploaded successfully",
		FileName: filename,
		Path:     dst,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})

}

func (fc *FileController) DownloadFile(c *gin.Context) {
	// Implement file download logic here
	fc.downloadMutex.Lock()
	defer fc.downloadMutex.Unlock()

	fileName := c.Query("file_name")
	directoryName := c.Query("directory_name")

	if fileName == "" || directoryName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	filePath := filepath.Join(directoryName, fileName)

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Type", "application/octet-stream")

	done := make(chan bool)
	errChan := make(chan error)

	go func() {
		file, err := os.Open(absPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

		buffer := make([]byte, 32*1024)
		for {
			n, err := file.Read(buffer)
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				errChan <- fmt.Errorf("failed to read file: %w", err)
				return
			}

			if _, err := c.Writer.Write(buffer[:n]); err != nil {
				errChan <- fmt.Errorf("failed to write file: %w", err)
				return
			}
			c.Writer.Flush()
		}
	}()

	select {
	case <-done:
		return
	case err := <-errChan:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	case <-time.After(5 * time.Minute):
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File download timed out"})
		return
	}

}
