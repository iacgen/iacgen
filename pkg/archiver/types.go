package archiver

type Archiver interface {
	Compress(filePath, pathToCompress string) error
	Decompress(filePath, pathToDecompress string) error
	Extension() string
}
