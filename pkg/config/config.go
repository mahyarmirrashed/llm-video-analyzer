package config

type Config struct {
	SamplingInterval int
	SamplingModel    string
	EmbeddingModel   string
	OllamaURL        string
	DatabaseURL      string
	Debug            bool
}
