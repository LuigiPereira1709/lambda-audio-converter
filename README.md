# Lambda Audio Converter
[Versão em Português](doc/README_pt.md)  
This project is a serverless audio converter built with AWS Lambda and Go. It uses FFmpeg to convert audio files to different formats and store them in MongoDB and S3. The project is designed to be cost-effective and efficient, leveraging the power of AWS Lambda and Go's performance. 

## Features
- [x] **Event-Driven**: Is triggered by S3 events, specifically when a metadata.json file is uploaded to the S3 bucket.
- [x] **Audio Conversion**: Converts audio files to different formats using FFmpeg.
- [x] **MongoDB Integration**: Update the document in MongoDB with the conversion results, including the duration of the audio file and the converted file's S3 URL. 
- [x] **S3 Integration**: Gets the files to convert from S3, delete the old files after conversion and store the converted files in S3.
- [x] **Metadata Handling**: Reads metadata from a JSON file uploaded to S3 and uses it to process the audio files.
- [x] **FFmpeg and FFprobe layer's**: Uses FFmpeg and FFprobe layers to handle audio processing efficiently.

## Workflow
1. **Trigger**: Triggered by an S3 event when a metadata.json file is uploaded to the S3 bucket.
2. **Event Parsing**: Parsing the S3 event to get the bucket name and necessary object keys.
3. **Metadata Retrieval**: Retrieves the metadata from the metadata.json file uploaded to S3.
4. **Get Objects**: Get the objects from event parsed data.
5. **Get Duration**: Uses FFprobe to get the duration of the audio file.
6. **Process Audio**: Converts the audio file to the desired format using FFmpeg.
7. **Delete Old Files**: Deletes the old audio files and metadata.json from S3 after conversion.
8. **Store Converted Files**: Stores the converted audio files in S3.
9. **Update Document**: Updates the MongoDB document based on the conversion success or failure. 
10. ** Cleanup**: Cleans up temporary files created during the process.

## Lambda Environment Variables
Example:
```json
{
  "Variables": {
    "MONGO_URI": "your_mongo_uri_with_credentials",
    "MONGO_DB": "your_database_name",

    "WORK_DIR": "/tmp/audio_converter",

    "FFMPEG_BIN_PATH": "/opt/bin/ffmpeg",
    "FFPROBE_BIN_PATH": "/opt/bin/ffprobe",

    "AUDIO_CONTENT_TYPE": "audio/m4a",
    "AUDIO_CODEC": "aac",
    "AUDIO_FORMAT": "m4a",

    "CONTENT_SUFFIX": ".m4a",
    "THUMBNAIL_SUFFIX": "thumbnail"
  }
}
```

## Notes:
The metadata.json file should be structured as follows:
```json
{
    "id": "unique_id",
    "title": "Audio Title",
    "year": "2003",
    "type": "music or podcast",
    "collection_name": "Collection Name",

    "music metadata": "below fields should be used only if type is music",
    "artist": "Artist Name",
    "album": "Album Name",
    "genre": "Genre",

    "podcast metadata": "below fields should be used only if type is podcast",
    "presenter": "Presenter Name",
    "description": "Podcast Description"
}
```

Organize your files in the S3 bucket as follows:
```plaintext
my-bucket/
└── document_id/
    └── document_title/
        ├── metadata.json   # Metadata file, this file will trigger the Lambda function, upload it last.
        ├── title.m4a     # Audio file converted to m4a format.
        ├── content.*       # Audio file in original format, include the extension, e.g., content.mp3.
        └── thumbnail       # Thumbnail file, omit the extension.
```

Log events will be generated in the CloudWatch logs, they will be similar to the following: 
```csv
`timestamp,message
1751235598686,"INIT_START Runtime Version: provided:al2023.v98	Runtime Version ARN: lambda-arn 
"
1751235598769,"START RequestId: <requestId> Version: $LATEST
"
1751235598978,"2025/06/29 22:19:58 Parsed event: {Bucket:<bucket name> ParentDirKey:<document id> EventFileKey:<document id>/metadata.json OthersFilesKey:map[content:<document id>/content thumbnail:<document id>/thumbnail]}
"
1751235598978,"2025/06/29 22:19:58 WORK_DIR doesn't exist, trying to create: /tmp/audio_converter
"
1751235599085,"2025/06/29 22:19:59 Parsed metadata: map[album:Album artist:Artist collection_name:music genre:ROCK id:<document_id> title:Music type:music year:2003]
"
1751235599710,"2025/06/29 22:19:59 Duration of the audio file: 199.079184 seconds
"
1751235599710,"2025/06/29 22:19:59 FFmpeg command: [/opt/bin/ffmpeg -y -progress pipe:1 -nostats -i /tmp/audio_converter/content -i /tmp/audio_converter/thumbnail -vf scale=trunc(iw/2)*2:trunc(ih/2)*2 -map 0:a -map 1:v -c:a aac -metadata:s:v title=Album cover -metadata:s:v comment=Cover (front) -metadata title=Music -metadata year=2003 -metadata artist=Artist -metadata album=Album -metadata genre=ROCK -movflags faststart /tmp/audio_converter/processed_file.m4a]
"
1751235619341,"2025/06/29 22:20:19 INFO File processed successfully details=""Progress: 98.35%. Current Time: 00:03:15. Duration: 199.08s. Finished: true. Elapsed Time: 00:00:19. Current Line: progress=end""
"
1751235619479,"2025/06/29 22:20:19 Content uploaded successfully to S3: <bucket_name>/<document_id>/Music.m4a
"
1751235619606,"2025/06/29 22:20:19 INFO Connected to MongoDB successfully dbName=<database_name>
"
1751235619610,"2025/06/29 22:20:19 Document updated successfully: {ID:<document_id> CollectionName:music ContentKey:<document_id>/Music.m4a Duration:199.079184 Status:0}
"
1751235619612,"2025/06/29 22:20:19 Temporary files cleaned up successfully
"
1751235619612,"2025/06/29 22:20:19 INFO Lambda handler completed successfully
"
1751235619613,"END RequestId: <requestId> 
"
1751235619613,"REPORT RequestId: <requestId>	Duration: 20844.05 ms	Billed Duration: 20925 ms	Memory Size: 640 MB	Max Memory Used: 113 MB	Init Duration: 80.95 ms	
"
```

## Project Structure
```plaintext
.
├── doc         # Extra documentation (Scripts)
├── handler     # Lambda function handler
├── internal
│   ├── converter    # Audio conversion and build logic
│   │   ├── music    # Music command build logic
│   │   └── podcast  # Podcast command build logic 
│   ├── database     # Database Connection
│   ├── s3      # S3 Service
│   └── utils        # Utility functions
├── main.go     # Main entry point for the Lambda function
└── scripts     # Scripts to build and deploy the Lambda function and FFmpeg layer
```

## Future Improvements
- [ ] Write unit and integration tests to ensure the functionality works as expected.
- [ ] Add support to convert audio files to other formats, such dash, opus, ogg, etc.
- [ ] Better goroutine management to handle multiple audio files concurrently.
- [ ] Integration with SNS and SQS for better event handling and error management.

## Links
- [AWS Lambda doc](https://aws.amazon.com/lambda/)
- [AWS S3 events doc](https://docs.aws.amazon.com/lambda/latest/dg/with-s3.html)
- [Why I choose Go for Lambda?](https://blog.scanner.dev/serverless-speed-rust-vs-go-java-python-in-aws-lambda-functions/)
- [FFmpeg doc](https://ffmpeg.org/ffmpeg.html)
- [FFprobe doc](https://ffmpeg.org/ffprobe.html)
- [Scripts doc](doc/scripts/scripts_doc_en.md)
- [Main doc](https://github.com/LuigiPereira1709/streaming-cloudnative-project)

## Learning Goals
- [x] Learn how to build a serverless audio converter using AWS Lambda and Go.
- [x] Understand how to use FFmpeg and FFprobe for audio processing.
- [x] Learn how to integrate AWS services like S3 and MongoDB with Go.
- [x] Gain experience in building and deploying serverless applications.
- [x] Learn how to handle S3 events and process audio files in a serverless environment.
- [x] Understand how to manage environment variables in AWS Lambda.
- [x] Learn how to use Lambda layers to include external libraries like FFmpeg and FFprobe.
- [ ] Learn how to optimize serverless applications for cost and performance.
- [ ] Gain experience in using Go's concurrency features to handle multiple audio files concurrently.
- [ ] Learn SQS and SNS integration for better event handling, notification and error management.
- [ ] Gain experience in AWS monitoring and logging services to track the performance and errors of the Lambda function.

## License
This project is licensed under the GNU GPL v3.0. See the [LICENSE](LICENSE.txt) file for details.
