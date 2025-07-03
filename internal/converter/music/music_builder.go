package music

import "pitanguinha.com/audio-converter/internal/converter"

// BuildCommand constructs the FFmpeg command for processing music files.
func BuildCommand(inputsPaths, metadataMap map[string]string) ([]string, error) {
	ffmpegCommand, err := converter.NewFFmpegCommand(inputsPaths, metadataMap, []string{"artist", "album", "genre"})
	if err != nil {
		return nil, err
	}

	ffmpegCommand.AddMetadataFromMap([]string{"artist", "album", "genre"}, metadataMap)

	return ffmpegCommand.BuildCommand(), nil
}
