package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"
)

func ShouldIgnoreFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	ignoredExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".mp4": true,
		".mov": true, ".mp3": true, ".webp": true, ".mkv": true, ".avi": true,
	}
	return ignoredExts[ext]
}

func CalculateSHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// EscapeMarkdown екранує спеціальні символи Telegram MarkdownV1: _ * ` [
func EscapeMarkdown(s string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"`", "\\`",
		"[", "\\[",
	)
	return replacer.Replace(s)
}
