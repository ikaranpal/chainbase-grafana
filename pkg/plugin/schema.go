package plugin

type Query struct {
	Statement string `json:"statement"`
}

type Options struct {
	QueriesPerSecond int    `json:"queriesPerSecond"`
	APIKey           string `json:"apiKey"`
}
