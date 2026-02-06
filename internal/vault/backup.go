package vault

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// CreateBackup creates a timestamped backup of the vault file.
func CreateBackup(vaultPath string, maxBackups int) error {
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return nil // Nothing to backup
	}

	timestamp := time.Now().Format("20060102150405")
	backupPath := fmt.Sprintf("%s.%s.bak", vaultPath, timestamp)

	if err := copyFile(vaultPath, backupPath); err != nil {
		return err
	}

	return cleanupOldBackups(vaultPath, maxBackups)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func cleanupOldBackups(vaultPath string, maxBackups int) error {
	dir := filepath.Dir(vaultPath)
	pattern := filepath.Base(vaultPath) + ".*.bak"
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return err
	}

	if len(matches) <= maxBackups {
		return nil
	}

	sort.Strings(matches)
	for i := 0; i < len(matches)-maxBackups; i++ {
		os.Remove(matches[i])
	}

	return nil
}
