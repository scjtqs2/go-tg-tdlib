package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	conf := DefaultConfig()
	if conf == nil {
		t.Fatal("expected non-nil default config")
	}
	if conf.AppID == "" {
		t.Error("expected non-empty AppID")
	}
	if conf.Proxy == nil {
		t.Error("expected non-nil Proxy config")
	}
	if conf.WebApi == nil {
		t.Error("expected non-nil WebApi config")
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "test_config.json")

	conf := DefaultConfig()
	conf.Phone = "+8618600000000"
	conf.Proxy.ProxyStatus = true
	conf.Proxy.ProxyType = "Socks5"

	err := conf.Save(confPath)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	loaded := Load(confPath)
	if loaded == nil {
		t.Fatalf("failed to load config from %s", confPath)
	}
	if loaded.Phone != "+8618600000000" {
		t.Errorf("expected phone '+8618600000000', got '%s'", loaded.Phone)
	}
	if !loaded.Proxy.ProxyStatus {
		t.Errorf("expected proxy status true")
	}
}

func TestCronJobConfigSave(t *testing.T) {
	tmpDir := t.TempDir()
	p := filepath.Join(tmpDir, "test_cron.json")

	cronCfg := &CronJobConfig{
		Cron: []*CronMessage{
			{
				Cron:       "0 * * * *",
				ToUserName: "@TestBot",
				TextMsg:    "/ping",
			},
		},
	}

	err := cronCfg.Save(p)
	if err != nil {
		t.Fatalf("failed to save cron config: %v", err)
	}

	if _, err := os.Stat(p); os.IsNotExist(err) {
		t.Fatalf("expected cron config file to exist at %s", p)
	}
}
