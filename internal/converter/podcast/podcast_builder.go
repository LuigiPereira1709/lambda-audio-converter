package podcast

import "pitanguinha.com/audio-converter/internal/converter"

// BuildCommand constructs the FFmpeg command for processing podcast files.
func BuildCommand(inputsPaths, metadataMap map[string]string) ([]string, error) {
	ffmpegCommand, err := converter.NewFFmpegCommand(inputsPaths, metadataMap, []string{"presenter", "description"})
	if err != nil {
		return nil, err
	}

	ffmpegCommand.AddMetadataFromMap([]string{"presenter", "description"}, metadataMap)

	return ffmpegCommand.BuildCommand(), nil
}
