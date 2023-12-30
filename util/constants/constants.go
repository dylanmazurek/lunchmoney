package constants

type ConstantEnum string
type ConstantString string
type ConstantInt int

var Config = newConfigRegistry()

type configRegistry struct {
	SourceUserAgent string
	APIBaseURL      string
}

func newConfigRegistry() *configRegistry {
	configRegistry := &configRegistry{
		"github.com/dylanmazurek/lunchmoney",
		"https://dev.lunchmoney.app/v1",
	}

	return configRegistry
}

var Path = newPathRegistry()

type pathRegistry struct {
	Me           string
	Assets       string
	Categories   string
	Transactions string
	Tags         string
}

func newPathRegistry() *pathRegistry {
	pathRegistry := &pathRegistry{
		"me",
		"assets",
		"categories",
		"transactions",
		"tags",
	}

	return pathRegistry
}

var ClientState = newClientStateRegistry()

type clientStateRegistry struct {
	New           string
	Initialized   string
	Authenticated string
	Error         string
}

func newClientStateRegistry() *clientStateRegistry {
	clientStateRegistry := &clientStateRegistry{
		"NEW",
		"INITIALIZED",
		"AUTHENTICATED",
		"ERROR",
	}

	return clientStateRegistry
}
