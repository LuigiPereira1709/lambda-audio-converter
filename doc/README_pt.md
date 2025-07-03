# Lambda Audio Converter
[English Version](../README.md)  
Este projeto é um conversor de áudio serverless construído com AWS Lambda e Go. Ele usa FFmpeg para converter arquivos de áudio em diferentes formatos e armazená-los no MongoDB e S3. O projeto foi projetado para ser econômico e eficiente, aproveitando o poder do AWS Lambda e o desempenho do Go.

## Funcionalidades
- [x] **Event-Driven**: É acionado por eventos do S3, especificamente quando um arquivo metadata.json é carregado no bucket S3.
- [x] **Conversão de Áudio**: Converte arquivos de áudio em diferentes formatos usando FFmpeg.
- [x] **Integração com MongoDB**: Atualiza o documento no MongoDB com os resultados da conversão, incluindo a duração do arquivo de áudio e a URL do arquivo convertido no S3.
- [x] **Integração com S3**: Obtém os arquivos a serem convertidos do S3, exclui os arquivos antigos após a conversão e armazena os arquivos convertidos no S3.
- [x] **Manipulação de Metadados**: Lê os metadados de um arquivo JSON carregado no S3 e os usa para processar os arquivos de áudio.
- [x] **FFmpeg e FFprobe layers**: Usa layers FFmpeg e FFprobe para lidar com o processamento de áudio de forma eficiente.

## Fluxo de Trabalho
1. **Gatilho**: Acionado por um evento do S3 quando um arquivo metadata.json é carregado no bucket S3.
2. **Parsing do Evento**: Parsing do evento S3 para obter o nome do bucket e as chaves dos objetos necessários. 
3. **Recuperação de Metadados**: Recupera os metadados do arquivo metadata.json carregado no S3.
4. **Get Objects**: Obtém os objetos a partir dos dados do evento analisado.
5. **Get Duration**: Usa FFprobe para obter a duração do arquivo de áudio.
6. **Processamento de Áudio**: Converte o arquivo de áudio para o formato desejado usando FFmpeg.
7. **Exclusão de Arquivos Antigos**: Exclui os arquivos de áudio antigos e o metadata.json do S3 após a conversão.
8. **Armazenamento de Arquivos Convertidos**: Armazena os arquivos de áudio convertidos no S3.
9. **Atualização de Documento**: Atualiza o documento do MongoDB com base no sucesso ou falha da conversão.
10. **Limpeza**: Limpa os arquivos temporários criados durante o processo.

## Variáveis de Ambiente do Lambda
Exemplo:
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

## Notas:
O arquivo metadata.json deve ser estruturado da seguinte forma:
```json
{
    "id": "unique_id",
    "title": "Titulo do Áudio",
    "year": "2003",
    "type": "music or podcast",
    "collection_name": "Nome da Coleção",

    "music metadata": "abaixo campos que devem ser usados apenas se o tipo for music",
    "artist": "Nome do Artista",
    "album": "Nome do Álbum",
    "genre": "Gênero da Música",

    "podcast metadata": "abaixo campos que devem ser usados apenas se o tipo for podcast",
    "presenter": "Nome do Apresentador",
    "description": "Descrição do Podcast"
}
```

Organize seus arquivos no bucket S3 da seguinte forma:
```plaintext
my-bucket/
└── document_id/
    └── document_title/
        ├── metadata.json       # Arquivo de metadados, esse arquivo ira disparar o lambda, ele deve ser o ultimo a ser carregado
        ├── title.m4a           # Arquivo de áudio convertido para o formato m4a.
        ├── content.*           # Arquivo de áudio no formato original, inclua a extensão, ex: content.mp3.
        └── thumbnail           # Arquivo de thumbnail, não inclua a extensão.
```

Os logs do evento serão gerados no CloudWatch, eles serão semelhantes ao seguinte:
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

## Estrutura do Projeto
```plaintext
.
├── doc         # Documentação extra (Scripts) 
├── handler     # Função Lambda handler 
├── internal
│   ├── converter    # Lógica de conversão de áudio e build de comandos FFmpeg
│   │   ├── music    # Lógica de build de commandos para music
│   │   └── podcast  # Lógica de build de commandos para podcast 
│   ├── database     # Conexão com o banco de dados 
│   ├── s3      # S3 Service
│   └── utils        # Funções utilitárias 
├── main.go     # Ponto de entrada principal para a função Lambda 
└── scripts     # Scripts para build e deploy da função Lambda e do layer FFmpeg 
```

## Melhorias Futuras
- [ ] Escrever testes unitários e de integração para a função Lambda.
- [ ] Adicionar suporte para mais formatos de áudio além de m4a.
- [ ] Melhor gerenciamento de goroutinas para processamento paralelo de arquivos.
- [ ] Integração com outros serviços AWS, como SNS ou SQS, para notificações e gerenciamento de erros. 

## Links
- [AWS Lambda doc](https://aws.amazon.com/lambda/)
- [AWS S3 events doc](https://docs.aws.amazon.com/lambda/latest/dg/with-s3.html)
- [Por que eu escolhi Go para a Lambda?](https://blog.scanner.dev/serverless-speed-rust-vs-go-java-python-in-aws-lambda-functions/)
- [FFmpeg doc](https://ffmpeg.org/ffmpeg.html)
- [FFprobe doc](https://ffmpeg.org/ffprobe.html)
- [Scripts doc](scripts/scripts_doc_pt.md)
- [Main doc](https://github.com/LuigiPereira1709/streaming-cloudnative-project/blob/main/doc/README_pt.md)

## Objetivos de Aprendizado
- [x] Aprender a construir uma função Lambda em Go que converte arquivos de áudio.
- [x] Compreender como usar o FFmpeg para processamento de áudio.
- [x] Aprender a integrar AWS Lambda com S3 e MongoDB.
- [x] Ganhar experiência prática com o desenvolvimento de aplicações serverless.
- [x] Aprender como lidar com eventos do S3 e processar arquivos de áudio de forma eficiente.
- [x] Compreender como usar variáveis de ambiente em funções Lambda para configurar o comportamento da aplicação.
- [x] Aprender como se usa layers no AWS Lambda para incluir dependências externas como FFmpeg e FFprobe.
- [ ] Aprender como otimizar aplicações serverless para reduzir custos e melhorar o desempenho.
- [ ] Ganhar experiência usando as funcionalidades concorrentes do Go para lidar com múltiplos arquivos de áudio simultaneamente.
- [ ] Aprender SQS e SNS para gerenciar eventos e notificações em aplicações serverless.
- [ ] Ganhar experiência com os serviços de monitoramento e logging do AWS Lambda para depuração e análise de desempenho. 

## Licença
Este projeto é licenciado sob a Licença GNU GPL v3.0. Veja o arquivo [LICENSE](../LICENSE.txt) para mais detalhes. 
