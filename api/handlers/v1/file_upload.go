package v1

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"strings"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/helpers"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/models"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
	static "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/statics"

	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"bytes"
	"image"

	_ "image/jpeg"
	_ "image/png"

	"github.com/nfnt/resize"

	webp "github.com/chai2010/webp"
)

const (
	maxFileSize    = 500 * 1024 // Target compressed file size (500KB)
	initialQuality = 50         // Starting WebP quality
	minQuality     = 10         // Lowest WebP quality allowed
	scaleFactor    = 0.5        // Default downscale (50%)
)

// @Router /v1/file/upload [post]
// @Summary Upload File to MinIO
// @Description API for uploading a file to MinIO and returning its URL
// @Tags upload
// @Security BearerAuth
// @Accept  multipart/form-data
// @Produce  json
// @Param file formData file true "File to upload"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "Invalid file"
// @Failure 500 {object} models.Response "Failed to upload file"
func (h *handlerV1) UploadFile(c *gin.Context) {
	h.log.Info("Uploading image...")

	// Initialize MinIO Client
	minioClient, err := minio.New(static.MiniOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(static.MiniOAccessKey, static.MiniOSecretKey, ""),
		Secure: static.MiniOProtocol,
	})
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "MinIO client initialization failed") {
		return
	}

	// Parse File from Form Data
	file, err := c.FormFile("file")
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Failed to read uploaded file") {
		return
	}

	// Open File Stream
	object, err := file.Open()
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error opening file") {
		return
	}
	defer object.Close()

	// Generate Unique File Name
	fName, _ := uuid.NewRandom()
	file.Filename = fmt.Sprintf("%s%s", fName.String(), filepath.Ext(file.Filename))
	file.Filename = strings.ReplaceAll(file.Filename, " ", "")

	// Read File into Buffer
	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, object); err != nil {
		helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error reading file into buffer")
		return
	}

	// **Compress Image Before Uploading**
	compressedImage, err := compressImage(buffer.Bytes()) // Start at 50% quality
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error compressing image") {
		return
	}

	// Upload to MinIO
	reader := bytes.NewReader(compressedImage)
	_, err = minioClient.PutObject(
		context.Background(),
		static.MiniOBucket,
		file.Filename,
		reader,
		reader.Size(),
		minio.PutObjectOptions{ContentType: "image/webp"},
	)
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Failed to upload file to MinIO") {
		return
	}

	// Response
	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "File uploaded successfully",
		Data: models.FileUploadReponse{
			Storage:          static.MiniOBucket,
			FileNameDisk:     file.Filename,
			FileNameDownload: file.Filename,
			Link:             fmt.Sprintf("%s/%s/%s", static.MiniOEndpoint, static.MiniOBucket, file.Filename),
			FileSize:         int64(len(compressedImage)),
		},
	})
}

// **Compress & Resize Image**
func compressImage(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	img, format, err := image.Decode(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// Resize if larger than 1000x1000
	if img.Bounds().Dx() > 1000 || img.Bounds().Dy() > 1000 {
		scale := 0.5
		newWidth := int(float64(img.Bounds().Dx()) * scale)
		newHeight := int(float64(img.Bounds().Dy()) * scale)
		img = resize.Resize(uint(newWidth), uint(newHeight), img, resize.Bilinear)
	}

	// Convert to WebP with optimized quality
	if format == "png" {
		err = webp.Encode(&buf, img, &webp.Options{Lossless: true}) // ✅ Lossless for PNGs
	} else {
		err = webp.Encode(&buf, img, &webp.Options{Quality: 35}) // ✅ Compressed for JPEGs
	}
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// @Router /v1/file/delete [delete]
// @Summary Delete File from MinIO
// @Description API for deleting a file from MinIO
// @Tags upload
// @Produce  json
// @Security BearerAuth
// @Param filename query string true "Filename to delete"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "Invalid request"
// @Failure 500 {object} models.Response "Failed to delete file"
func (h *handlerV1) DeleteFile(c *gin.Context) {
	filename := strings.Replace(c.Query("filename"), "minio.taklifnomavip.uz/my-bucket/", "", 1)
	if filename == "" {
		c.JSON(400, models.Response{
			Code:    config.ErrorBadRequest,
			Message: "Filename is required",
		})
		return
	}

	minioClient, err := minio.New(static.MiniOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(static.MiniOAccessKey, static.MiniOSecretKey, ""),
		Secure: static.MiniOProtocol,
	})
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Failed to create minio client") {
		return
	}

	_, err = minioClient.StatObject(context.Background(), static.MiniOBucket, filename, minio.StatObjectOptions{})
	if err != nil {
		c.JSON(400, models.Response{
			Code:    config.ErrorBadRequest,
			Message: fmt.Sprintf("File not found in MinIO: %v", err),
		})
		return
	}

	err = minioClient.RemoveObject(
		context.Background(),
		static.MiniOBucket,
		strings.Replace(filename, "minio.taklifnomavip.uz/my-bucket/", "", 1),
		minio.RemoveObjectOptions{
			GovernanceBypass: true,
		},
	)
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Failed to delete file from minio") {
		return
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "File deleted successfully",
		Data: map[string]string{
			"deleted_file": filename,
		},
	})
}
