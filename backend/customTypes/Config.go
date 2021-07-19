package customTypes

type Config struct {
	Node struct {
		Endpoint      string `yaml:"endpoint" env:"TQ_NODE_ENDPOINT" env-description:"Endpoint on where the Filecoin node is hosted"`
		TokenRequired bool   `yaml:"token-required" env:"TQ_NODE_TOKEN_REQUIRED" env-description:"Is token required to connect ?"`
		Token         string `yaml:"token" env:"TQ_NODE_TOKEN" env-description:"Token to connect to Filecoin node"`
	} `yaml:"node"`

	Database struct {
		EndpointWithCredentials string `yaml:"endpoint-with-creds" env:"TQ_DB_ENDPOINT" env-description:"Database endpoint"`
		DatabaseName            string `yaml:"database-name" env:"TQ_DB_NAME" env-description:"Database name for Slack bot"`
	} `yaml:"database"`

	Slack struct {
		BotToken      string `yaml:"bot-token" env:"TQ_SLACK_BOT_TOKEN" env-description:"Slack bot token to connect to Slack API"`
		SigningSecret string `yaml:"signing-secret" env:"TQ_SLACK_SIGNING_SECRET" env-description:"Slack bot signing secret for auth"`
	} `yaml:"slack"`
}
