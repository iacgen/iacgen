package model

type ProjectDetails struct {
	Metadata RepoMetadata        `json:"metadata"`
	Graph    map[string][]string `json:"graph"`
	Projects []Project           `json:"projects"`
}
type RepoMetadata struct {
	Repository string `json:"repository"`
	Commit     string `json:"commit"`
	Branch     string `json:"branch"`
}
type Metadata struct {
	Location   string   `json:"location"`
	Name       string   `json:"name"`
	Languages  []string `json:"languages"`
	Frameworks []string `json:"frameworks"`
}
type Ports struct {
	Listen  int `json:"listen"`
	Forward int `json:"forward,omitempty"`
}
type HTTPEgress struct {
	Endpoint   string   `json:"endpoint"`
	Operations []string `json:"operations"`
}
type HTTPIngress struct {
	Endpoint   string   `json:"endpoint"`
	Operations []string `json:"operations"`
}
type S3 struct {
	Bucket     string   `json:"bucket"`
	Operations []string `json:"operations"`
}
type Database struct {
	Dsn string `json:"dsn"`
}
type Services struct {
	Name        string        `json:"name"`
	DNS         string        `json:"dns"`
	Image       string        `json:"image,omitempty"`
	Ports       []Ports       `json:"ports"`
	HTTPEgress  []HTTPEgress  `json:"http_egress"`
	HTTPIngress []HTTPIngress `json:"http_ingress,omitempty"`
	S3          []S3          `json:"s3"`
	Database    []Database    `json:"database"`
}
type Project struct {
	Metadata Metadata   `json:"metadata"`
	Services []Services `json:"services"`
}
