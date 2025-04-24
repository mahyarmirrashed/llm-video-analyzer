package config

type Config struct {
	SamplingInterval int
	SamplingModel    string
	EmbeddingModel   string
	QueryLimit       int
	QueryModel       string
	OllamaURL        string
	DatabaseURL      string
	Debug            bool
}
