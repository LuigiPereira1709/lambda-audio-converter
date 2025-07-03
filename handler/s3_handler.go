package handler

import (
	"pitanguinha.com/audio-converter/internal/s3"
	"pitanguinha.com/audio-converter/internal/utils"
)

// GetFilesFromS3 retrieves the metadata, thumbnail, and content files from S3 based on the event parsed.
func GetFilesFromS3(s3Service *s3.S3Service, eventParsed EventParsed) (map[string]string, error) {
	type fileSpec struct {
		s3KeyName string
		fileName  string
	}

	fileSpecs := []fileSpec{
		{s3KeyName: eventParsed.EventFileKey, fileName: "metadata"},
		{s3KeyName: eventParsed.OthersFilesKey["thumbnail"], fileName: "thumbnail"},
		{s3KeyName: eventParsed.OthersFilesKey["content"], fileName: "content"},
	}

	filesPaths := make(map[string]string)
	workDir := utils.GetWorkDir()

	for _, spec := range fileSpecs {
		reader, err := s3Service.GetObject(eventParsed.Bucket, spec.s3KeyName)
		if err != nil {
			return nil, err
		}

		filePath, err := utils.WriteToFileFromReader(workDir, spec.fileName, reader)
		if err != nil {
			return nil, err
		}
		filesPaths[spec.fileName] = filePath
	}
	return filesPaths, nil
}

// UploadContentToS3 uploads the content file to the specified S3 bucket with the given key and content type.
func UploadContentToS3(s3Service *s3.S3Service, bucket, key, contentType, filePath string) error {
	file, err := utils.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return s3Service.PutObject(bucket, key, contentType, file)
}

// DeleteFilesFromS3 deletes the specified files from the S3 bucket.
func DeleteFilesFromS3(s3Service *s3.S3Service, bucket string, keys ...string) error {
	for _, key := range keys {
		if err := s3Service.DeleteObject(bucket, key); err != nil {
			return err
		}
	}
	return nil
}
