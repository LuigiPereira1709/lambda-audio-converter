package handler

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"pitanguinha.com/audio-converter/internal/converter"
	"pitanguinha.com/audio-converter/internal/s3"
	"pitanguinha.com/audio-converter/internal/utils"
)

// parseMetadata reads and parses the metadata file from S3.
func parseMetadata(metadataPath string) (map[string]string, error) {
	data, err := utils.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("error reading metadata file %s: %w", metadataPath, err)
	}

	rawMap, err := utils.ParseJsonToMap(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing metadata file %s: %w", metadataPath, err)
	}

	metadata := make(map[string]string)
	for k, v := range rawMap {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("metadata field %s is not a string", k)
		}
		metadata[k] = str
	}

	// Validate required fields
	required := []string{"id", "title", "collection_name"}
	for _, key := range required {
		if metadata[key] == "" {
			return nil, fmt.Errorf("missing required metadata field: %s", key)
		}
	}

	return metadata, nil
}

func encodeContentKey(contentKey string) string {
	lastSlash := strings.LastIndex(contentKey, "/")
	folder := contentKey[:lastSlash]
	fileName := contentKey[lastSlash+1:]
	return fmt.Sprintf("%s/%s", folder, url.PathEscape(fileName))
}

// Handler processes an audio conversion Lambda event.
func Handler(ctx context.Context, event events.S3Event) error {
	audioContentType := os.Getenv("AUDIO_CONTENT_TYPE")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		slog.Error("failed to load AWS config", "err", err)
		return nil
	}

	s3Service := s3.NewService(cfg)

	eventParsed, err := ParseEvent(s3Service, event)
	if err != nil {
		slog.Error("error parsing event", "err", err)
		return nil
	}
	log.Printf("Parsed event: %+v", eventParsed)

	filesPaths, err := GetFilesFromS3(s3Service, eventParsed)
	if err != nil {
		slog.Error("error getting files from S3", "err", err)
		return nil
	}

	metadata, err := parseMetadata(filesPaths["metadata"])
	if err != nil {
		slog.Error("error parsing metadata", "err", err)
		return nil
	}
	log.Printf("Parsed metadata: %+v", metadata)

	duration, err := converter.GetDurationFromFile(filesPaths["content"])
	if err != nil {
		slog.Error("error getting duration", "err", err)
		return nil
	}
	log.Printf("Duration of the audio file: %f seconds", duration)

	details, err := ProcessAudioFile(duration, filesPaths, metadata)
	if err != nil {
		slog.Error("error processing audio file", "err", err, "details", details)
		return nil
	}
	slog.Info("File processed successfully", "details", details)

	bucket := eventParsed.Bucket
	keysToDelete := []string{
		eventParsed.EventFileKey,
		eventParsed.OthersFilesKey["content"],
	}

	if err := DeleteFilesFromS3(s3Service, bucket, keysToDelete...); err != nil {
		slog.Error("error deleting old files from S3", "err", err)
		return nil
	}

	contentKey := fmt.Sprintf("%s/%s.%s", eventParsed.ParentDirKey, metadata["title"], os.Getenv("AUDIO_FORMAT"))
	if err := UploadContentToS3(s3Service, bucket, contentKey, audioContentType, details.ProcessedFilePath); err != nil {
		slog.Error("error uploading converted content to S3", "bucket", bucket, "key", contentKey, "err", err)
		return nil
	}
	log.Printf("Content uploaded successfully to S3: %s/%s", bucket, contentKey)

	doc := UpdateDocumentInput{
		ID:             metadata["id"],
		CollectionName: metadata["collection_name"],
		ContentKey:     encodeContentKey(contentKey),
		Duration:       duration,
		Status:         SetStatus(details.Finished),
	}

	if err := doc.UpdateDocument(); err != nil {
		slog.Error("error updating document in database", "err", err)
		return nil
	}
	log.Printf("Document updated successfully: %+v", doc)

	if err := utils.DeleteFiles(utils.GetWorkDir()); err != nil {
		slog.Warn("failed to clean up temporary files", "err", err)
	}
	log.Printf("Temporary files cleaned up successfully")

	slog.Info("Lambda handler completed successfully")
	return nil
}
