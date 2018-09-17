package config

import (
	"net/url"

	"os"

	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Config is the configuration for horizon.  It get's populated by the
// app's main function and is provided to NewApp.
type Config struct {
	*Base
	Hostname               string
	DatabaseURL            string
	StellarCoreDatabaseURL string
	StellarCoreURL         string
	Port                   int
	RedisURL               string
	LogLevel               logrus.Level
	LogToJSON              bool
	SlowQueryBound         *time.Duration

	APIBackend      *url.URL
	KeychainBackend *url.URL

	SlackWebhook *url.URL
	SlackChannel string

	// TLSCert is a path to a certificate file to use for horizon's TLS config
	TLSCert string
	// TLSKey is the path to a private key file to use for horizon's TLS config
	TLSKey string
	// Ingest is a boolean that indicates whether or not this horizon instance
	// should run the data ingestion subsystem.
	Ingest bool
	// HistoryRetentionCount represents the minimum number of ledgers worth of
	// history data to retain in the horizon database. For the purposes of
	// determining a "retention duration", each ledger roughly corresponds to 10
	// seconds of real time.
	HistoryRetentionCount uint
	// StaleThreshold represents the number of ledgers a history database may be
	// out-of-date by before horizon begins to respond with an error to history
	// requests.
	StaleThreshold uint
	//For developing without signatures
	SkipCheck bool
	// enable on dev only
	CORSEnabled bool
	// DisableAPISubmit tell horizon to not use API for transaction submission
	// for dev purposes only, works well with SkipCheck enabled
	// pending transactions and transaction 2fa will be disabled as well.
	DisableAPISubmit bool
	// If set to true - Horizon won't check TFA (via API) during TX submission.
	DisableTXTfa bool

	Core Core

	TemplateBackend *url.URL
	InvestReady     *url.URL
	TelegramAirdrop *url.URL

	ForceHTTPSLinks bool

	SentryDSN      string
	Project        string
	SentryLogLevel string
	Env            string

	MigrateUpOnStart bool
}

func (c *Config) DefineConfigStructure(cmd *cobra.Command) {
	c.Base = NewBase(nil, "")

	c.Core.Base = NewBase(c.Base, "core")
	c.Core.DefineConfigStructure()

	c.setDefault("port", 8000)
	c.setDefault("per_hour_hate_limit", 72000)
	c.setDefault("history_retention_count", 0)
	c.setDefault("sign_checkskip", false)
	c.setDefault("log_level", "debug")
	c.setDefault("force_https_links", true)
	c.setDefault("sentry_dsn", "")
	c.setDefault("project", "")
	c.setDefault("sentry_log_level", "warn")
	c.setDefault("env", "")

	c.bindEnv("port")
	c.bindEnv("database_url")
	c.bindEnv("api_database_url")
	c.bindEnv("stellar_core_database_url")
	c.bindEnv("stellar_core_url")
	c.bindEnv("per_hour_hate_limit")
	c.bindEnv("redis_url")
	c.bindEnv("log_level")
	c.bindEnv("log_to_json")
	c.bindEnv("slow_query_bound")

	c.bindEnv("tls_cert")
	c.bindEnv("tls_key")
	c.bindEnv("ingest")
	c.bindEnv("history_retention_count")
	c.bindEnv("history_stale_threshold")
	c.bindEnv("sign_check_skip")

	c.bindEnv("templates_path")

	c.bindEnv("horizon_secret")
	c.bindEnv("keyserver_url")

	c.bindEnv("api_backend")
	c.bindEnv("keychain_backend")

	c.bindEnv("slack_webhook")
	c.bindEnv("slack_channel")

	c.bindEnv("cors_enabled")
	c.bindEnv("hostname")

	c.bindEnv("disable_api_submit")

	c.bindEnv("template_backend")
	c.bindEnv("invest_ready")
	c.bindEnv("telegram_airdrop")
	c.bindEnv("disable_tx_tfa")
	c.bindEnv("force_https_links")
	c.bindEnv("migrate_up_on_start")
}

func (c *Config) Init() error {
	c.Port = c.getInt("port")

	var err error
	c.DatabaseURL, err = c.getNonEmptyString("database_url")
	if err != nil {
		return err
	}

	c.StellarCoreDatabaseURL, err = c.getNonEmptyString("stellar_core_database_url")
	if err != nil {
		return err
	}

	c.StellarCoreURL, err = c.getNonEmptyString("stellar_core_url")
	if err != nil {
		return err
	}

	c.RedisURL = c.getString("redis_url")

	c.SlowQueryBound, err = c.getOptionalTDuration("slow_query_bound")
	if err != nil {
		return err
	}

	c.LogToJSON = c.getBool("log_to_json")
	c.LogLevel, err = logrus.ParseLevel(c.getString("log_level"))
	if err != nil {
		return err
	}

	c.TLSCert = c.getString("tls_cert")
	c.TLSKey = c.getString("tls_key")
	switch {
	case c.TLSCert != "" && c.TLSKey == "":
		return errors.New("Invalid TLS config: key not configured")
	case c.TLSCert == "" && c.TLSKey != "":
		return errors.New("Invalid TLS config: cert not configured")
	}

	c.Ingest = c.getBool("ingest")
	c.HistoryRetentionCount = uint(c.getInt("history_retention_count"))
	c.StaleThreshold = uint(c.getInt("history_stale_threshold"))
	c.SkipCheck = c.getBool("sign_check_skip")

	c.RedisURL = c.getString("redis_url")

	err = c.Core.Init()
	if err != nil {
		return err
	}

	c.APIBackend, err = c.getParsedURL("api_backend")
	if err != nil {
		return err
	}

	c.KeychainBackend, err = c.getParsedURL("keychain_backend")
	if err != nil {
		return err
	}

	if c.getString("slack_webhook") != "" {
		c.SlackWebhook, err = c.getParsedURL("slack_webhook")
		if err != nil {
			return err
		}
		c.SlackChannel, err = c.getNonEmptyString("slack_channel")
		if err != nil {
			return err
		}
	}

	c.CORSEnabled = c.getBool("cors_enabled")
	c.Hostname = c.getString("hostname")
	if c.Hostname == "" {
		c.Hostname, err = os.Hostname()
		if err != nil {
			return errors.Wrap(err, "failed to get hostname")
		}
	}

	c.DisableAPISubmit = c.getBool("disable_api_submit")
	c.DisableTXTfa = c.getBool("disable_tx_tfa")

	c.TemplateBackend, err = c.getParsedURL("template_backend")
	if err != nil {
		return errors.Wrap(err, "Failed to get template_backend value")
	}

	c.InvestReady, err = c.getParsedURL("invest_ready")
	if err != nil {
		return errors.Wrap(err, "Failed to get invest_ready value")
	}

	c.TelegramAirdrop, err = c.getOptionalParsedURL("telegram_airdrop")
	if err != nil {
		return errors.Wrap(err, "Failed to get telegram_airdrop value")
	}

	c.ForceHTTPSLinks = c.getBool("force_https_links")
	c.SentryDSN = c.getString("sentry_dsn")
	c.Project = c.getString("project")
	c.SentryLogLevel = c.getString("sentry_log_level")

	c.MigrateUpOnStart = c.getBool("migrate_up_on_start")
	c.Env = c.getString("env")
	return nil
}
