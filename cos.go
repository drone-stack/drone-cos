package cos

type (
	Cos struct {
		Bucket      string
		AccessKey   string
		SecretKey   string
		Region      string
		Source      string
		Target      string
		StripPrefix string
		Endpoint    string
	}

	Plugin struct {
		Cos Cos
	}
)

// Exec executes the plugin step
func (p Plugin) Exec() error {
	return nil
}
