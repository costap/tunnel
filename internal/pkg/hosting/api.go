package hosting

type Provider interface {
	CreateInstance() (string, error)
}
