package archiver

type Archiver interface {
	Compress(filePath, pathToCompress string) error
	Extension() string
}
