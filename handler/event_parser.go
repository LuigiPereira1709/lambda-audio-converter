package handler

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"pitanguinha.com/audio-converter/internal/s3"
	"pitanguinha.com/audio-converter/internal/utils"
)

// EventParsed holds the parsed information from an S3 event.
type EventParsed struct {
	Bucket         string `json:"bucket.name"`
	ParentDirKey   string
	EventFileKey   string
	OthersFilesKey map[string]string
}

// ParseEvent parses the S3 event and retrieves the bucket name, event file key, and other files in the same directory.
func ParseEvent(s3Service *s3.S3Service, event events.S3Event) (EventParsed, error) {
	var eventParsed EventParsed

	eventParsed.Bucket = event.Records[0].S3.Bucket.Name

	decodeKey, err := url.QueryUnescape(event.Records[0].S3.Object.Key)
	if err != nil {
		return eventParsed, fmt.Errorf("error decoding S3 object key: %w", err)
	}
	eventParsed.EventFileKey = decodeKey

	dir := utils.GetParentDir(decodeKey)
	eventParsed.ParentDirKey = dir

	if err := eventParsed.loadAdditionalFileKeys(s3Service, dir); err != nil {
		return eventParsed, fmt.Errorf("error loading additional file keys: %w", err)
	}

	return eventParsed, nil
}

// loadAdditionalFileKeys retrieves the paths of other files in the same directory as the event file.
func (e *EventParsed) loadAdditionalFileKeys(s3Service *s3.S3Service, dir string) error {
	keys, err := s3Service.ListObjectsForPrefix(e.Bucket, dir) // NOTE: Expect: Event file, thumbnail file and one or two content files.
	if err != nil {
		return err
	}

	// INFO: On creation, we have 3 files (event file, content file and thumbnail file).
	// On update, we have 3 (same the creation) or 4 files (same the creation + content file wich will replace the old one).
	ContentSuffix := os.Getenv("CONTENT_SUFFIX")
	ThumbnailSuffix := os.Getenv("THUMBNAIL_SUFFIX")
	e.OthersFilesKey = make(map[string]string)
	var contentWithoutSuffix string

	for _, key := range keys {
		// Skip the event file itself and the content file if it has a specific suffix.
		if key == e.EventFileKey || key == e.ParentDirKey+"/" {
			continue
		}

		// Remove the spaces from the key
		key = strings.TrimSpace(key)

		if strings.HasSuffix(key, ThumbnailSuffix) {
			e.OthersFilesKey["thumbnail"] = key
			continue
		}

		if strings.HasSuffix(key, ContentSuffix) {
			e.OthersFilesKey["content"] = key
			continue
		}

		contentWithoutSuffix = key
	}

	if contentWithoutSuffix != "" {
		e.OthersFilesKey["content"] = contentWithoutSuffix
	}

	return nil
}
