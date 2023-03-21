package cors

type Cors struct {
	Config *Config
}

type Config struct {
	Origins []string
	Methods []string
}

func New(config Config) *Cors {
	return &Cors{
		Config: &config,
	}
}

func (cors *Cors) ValidateOrigin(requestOrigin string) string {

	for _, origin := range cors.Config.Origins {

		if origin == "*" {
			return "*"
		}

		if requestOrigin == origin {
			return origin
		}

	}

	return ""

}

// Return the list of allowed methods in a unique string
func (cors *Cors) AllowedMethods() string {

	var methods string

	for i, method := range cors.Config.Methods {
		methods += method

		if i < len(cors.Config.Methods)-1 {
			methods += ","
		}
	}

	return methods
}
