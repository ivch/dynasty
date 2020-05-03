package requests

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"

	"github.com/ivch/dynasty/common"
	"github.com/ivch/dynasty/models/dto"
)

const (
	allowedFileType = "image/jpeg"
	filesPerRequest = 3
	imgPathPrefix   = "req/i/"
	thumbPathPrefix = "req/t/"
	defaultS3ACL    = "public-read"
)

func (s *service) UploadImage(_ context.Context, r *dto.UploadImageRequest) (*dto.UploadImageResponse, error) {
	req, err := s.repo.GetRequestByIDAndUser(r.RequestID, r.UserID)
	if err != nil {
		return nil, err
	}

	if len(req.Images) >= filesPerRequest {
		return nil, errTooMuchFiles
	}

	filename := fmt.Sprintf("%s:%s.jpg", base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(r.UserID))), common.RandomString(25))
	imgPath := s.buildImagePath(imgPathPrefix, filename)
	imgThumbPath := s.buildImagePath(thumbPathPrefix, filename)

	fileType := http.DetectContentType(r.File)
	if fileType != allowedFileType {
		return nil, errFileWrongType
	}

	thumb, err := s.createThumbnail(r.File)
	if err != nil {
		return nil, err
	}

	if _, err := s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.s3Space),
		Key:         aws.String(imgPath),
		Body:        bytes.NewReader(r.File),
		ACL:         aws.String(defaultS3ACL),
		ContentType: aws.String(fileType),
	}); err != nil {
		return nil, err
	}

	if _, err := s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.s3Space),
		Key:         aws.String(imgThumbPath),
		Body:        bytes.NewReader(thumb),
		ACL:         aws.String(defaultS3ACL),
		ContentType: aws.String(fileType),
	}); err != nil {
		// todo this can cause an error because deleteImageFromS3 expects both img and thumb are present
		if s3err := s.deleteImageFromS3(filename); s3err != nil {
			return nil, s3err
		}
		return nil, err
	}

	if err := s.repo.AddImage(r.UserID, r.RequestID, filename); err != nil {
		if s3err := s.deleteImageFromS3(imgPath); s3err != nil {
			return nil, s3err
		}
		return nil, err
	}

	imgUrl := s.buildImageURL(filename)
	return &dto.UploadImageResponse{
		Img:   imgUrl["img"],
		Thumb: imgUrl["thumb"],
	}, nil
}

func (s *service) DeleteImage(_ context.Context, r *dto.DeleteImageRequest) error {
	filename := filepath.Base(r.Filepath)

	if err := s.repo.DeleteImage(r.UserID, r.RequestID, filename); err != nil {
		return err
	}

	if err := s.deleteImageFromS3(filename); err != nil {
		if err2 := s.repo.AddImage(r.UserID, r.RequestID, filename); err2 != nil {
			return err2
		}
		return err
	}

	return nil
}

func (s *service) deleteImageFromS3(filename string) error {
	img := s.buildImagePath(imgPathPrefix, filename)
	if _, err := s.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.s3Space),
		Key:    aws.String(img),
	}); err != nil {
		return fmt.Errorf("failed to delete image %s: %w", img, err)
	}

	thumb := s.buildImagePath(thumbPathPrefix, filename)
	if _, err := s.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.s3Space),
		Key:    aws.String(thumb),
	}); err != nil {
		return fmt.Errorf("failed to delete image %s: %w", thumb, err)
	}

	return nil
}

func (s *service) createThumbnail(file []byte) ([]byte, error) {
	img, err := imaging.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}

	var (
		thumb = imaging.Thumbnail(img, 128, 128, imaging.CatmullRom)
		dst   = imaging.New(128, 128, color.NRGBA{R: 0, G: 0, B: 0, A: 0})
		buf   = new(bytes.Buffer)
	)

	dst = imaging.Paste(dst, thumb, image.Pt(0, 0))
	if err := jpeg.Encode(buf, dst, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *service) buildImageURL(filename string) map[string]string {
	return map[string]string{
		"img":   fmt.Sprintf("%s/%s", s.cdnHost, s.buildImagePath(imgPathPrefix, filename)),
		"thumb": fmt.Sprintf("%s/%s", s.cdnHost, s.buildImagePath(thumbPathPrefix, filename)),
	}
}

func (s *service) buildImagePath(prefix, filename string) string {
	return fmt.Sprintf("%s%s", prefix, filename)
}
