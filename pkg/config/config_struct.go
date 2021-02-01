package config

import (
	"github.com/hiromaily/go-book-teacher/pkg/notifier"
	"github.com/hiromaily/go-book-teacher/pkg/save"
	"github.com/hiromaily/go-book-teacher/pkg/site"
)

// Root is root config
type Root struct {
	Interval     int           `toml:"interval" validate:"required"`
	Logger       *Logger       `toml:"logger" validate:"required"`
	Site         *Site         `toml:"site" validate:"required"`
	Save         *Save         `toml:"save"`
	Notification *Notification `toml:"notification"`
}

// Logger is zap logger property
type Logger struct {
	Service      string `toml:"service" validate:"required"`
	Env          string `toml:"env" validate:"oneof=dev prod custom"`
	Level        string `toml:"level" validate:"required"`
	IsStackTrace bool   `toml:"is_stacktrace"`
}

// Site is site information
type Site struct {
	Type site.SiteType `toml:"type" validate:"oneof=dmm"`
	URL  string        `toml:"url" validate:"required"`
}

// Save is save method
type Save struct {
	Mode  save.Mode `toml:"mode" validate:"oneof=text redis dummy"`
	Text  *Text     `toml:"text" validate:"-"`
	Redis *Redis    `toml:"redis" validate:"-"`
}

// Text is text save
type Text struct {
	Path string `toml:"path" validate:"required"`
}

// Redis is redis save
type Redis struct {
	Encrypted bool   `toml:"encrypted"`
	URL       string `toml:"url" validate:"required"`
}

// Notification is notification method
type Notification struct {
	Mode    notifier.Mode `toml:"mode" validate:"required"`
	Console *Console      `toml:"console" validate:"-"`
	Slack   *Slack        `toml:"slack" validate:"-"`
}

// Console is command line notification
type Console struct {
	Enabled bool `toml:"enabled"`
}

// Slack is slack notification
type Slack struct {
	Enabled   bool   `toml:"enabled"`
	Encrypted bool   `toml:"encrypted"`
	Key       string `toml:"key" validate:"required"`
}
