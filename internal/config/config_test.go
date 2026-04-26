package config

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestLoad_Success(t *testing.T) {
	t.Setenv("TELEGRAM_TOKEN", "test-token-123")
	t.Setenv("VT_API_KEY", "test-api-key-456")
	t.Setenv("VERBOSE_OUTPUT", "")

	cfg := Load()

	if cfg.TelegramToken != "test-token-123" {
		t.Errorf("TelegramToken = %q, want %q", cfg.TelegramToken, "test-token-123")
	}
	if cfg.VTApiKey != "test-api-key-456" {
		t.Errorf("VTApiKey = %q, want %q", cfg.VTApiKey, "test-api-key-456")
	}
	if cfg.MaxFileSize != 20*1024*1024 {
		t.Errorf("MaxFileSize = %d, want %d", cfg.MaxFileSize, 20*1024*1024)
	}
	if cfg.RequestInterval != 16*time.Second {
		t.Errorf("RequestInterval = %v, want %v", cfg.RequestInterval, 16*time.Second)
	}
	if cfg.VerboseOutput != false {
		t.Errorf("VerboseOutput = %v, want false", cfg.VerboseOutput)
	}
}

func TestLoad_VerboseOutputTrue(t *testing.T) {
	t.Setenv("TELEGRAM_TOKEN", "token")
	t.Setenv("VT_API_KEY", "key")
	t.Setenv("VERBOSE_OUTPUT", "true")

	cfg := Load()

	if cfg.VerboseOutput != true {
		t.Errorf("VerboseOutput = %v, want true", cfg.VerboseOutput)
	}
}

func TestLoad_VerboseOutputFalseVariants(t *testing.T) {
	variants := []string{"false", "1", "yes", "TRUE", "True", ""}

	for _, v := range variants {
		t.Run("VERBOSE_OUTPUT="+v, func(t *testing.T) {
			t.Setenv("TELEGRAM_TOKEN", "token")
			t.Setenv("VT_API_KEY", "key")
			t.Setenv("VERBOSE_OUTPUT", v)

			cfg := Load()

			// Тільки точне "true" повинно давати true
			if v == "true" && !cfg.VerboseOutput {
				t.Errorf("VerboseOutput = false, want true for input %q", v)
			}
			if v != "true" && cfg.VerboseOutput {
				t.Errorf("VerboseOutput = true, want false for input %q", v)
			}
		})
	}
}

// TestLoad_MissingToken перевіряє, що Load() завершує процес при відсутньому TELEGRAM_TOKEN.
// Використовуємо subprocess-підхід, бо log.Fatal викликає os.Exit(1).
func TestLoad_MissingToken(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		os.Setenv("TELEGRAM_TOKEN", "")
		os.Setenv("VT_API_KEY", "some-key")
		Load()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLoad_MissingToken")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=1", "TELEGRAM_TOKEN=", "VT_API_KEY=some-key")
	err := cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok && !exitErr.Success() {
		return // Очікуваний вихід з помилкою
	}
	t.Fatal("Load() повинен завершити процес при відсутньому TELEGRAM_TOKEN")
}

// TestLoad_MissingAPIKey перевіряє, що Load() завершує процес при відсутньому VT_API_KEY.
func TestLoad_MissingAPIKey(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		os.Setenv("TELEGRAM_TOKEN", "some-token")
		os.Setenv("VT_API_KEY", "")
		Load()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLoad_MissingAPIKey")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=1", "TELEGRAM_TOKEN=some-token", "VT_API_KEY=")
	err := cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok && !exitErr.Success() {
		return // Очікуваний вихід з помилкою
	}
	t.Fatal("Load() повинен завершити процес при відсутньому VT_API_KEY")
}

// TestLoad_BothMissing перевіряє, що Load() завершує процес, коли обидва ключі відсутні.
func TestLoad_BothMissing(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		os.Setenv("TELEGRAM_TOKEN", "")
		os.Setenv("VT_API_KEY", "")
		Load()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLoad_BothMissing")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=1", "TELEGRAM_TOKEN=", "VT_API_KEY=")
	err := cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok && !exitErr.Success() {
		return // Очікуваний вихід з помилкою
	}
	t.Fatal("Load() повинен завершити процес при відсутніх обох ключах")
}
