package commands

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/zulfikawr/gotp/internal/cli/ui"
	"github.com/zulfikawr/gotp/internal/config"
	"github.com/zulfikawr/gotp/internal/vault"
)

// executeCommand is a helper to run a cobra command and return all output.
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()

	uiOut := ""
	if b, ok := ui.Out.(*bytes.Buffer); ok {
		uiOut = b.String()
	}

	res := uiOut + buf.String()
	return res, err
}

func setupTestCLI(vaultPath string, input string) *cobra.Command {
	config.SetVaultPathOverride(vaultPath)
	ui.In = strings.NewReader(input)
	ui.Out = new(bytes.Buffer)
	
	// Mock terminal for tests
	ui.IsTerminal = func(fd int) bool { return true }
	ui.PasswordReader = func(fd int) ([]byte, error) {
		s := bufio.NewScanner(ui.In)
		if s.Scan() {
			val := s.Text()
			// recreate reader with remaining data for subsequent prompts
			remaining := ""
			for s.Scan() {
				remaining += s.Text() + "\n"
			}
			ui.In = strings.NewReader(remaining)
			return []byte(val), nil
		}
		return nil, io.EOF
	}
	ui.ResetScanner()
	ui.ResetScanner()

	// Clear session for each test to ensure predictable prompts
	_ = vault.ClearSession()

	root := &cobra.Command{Use: "gotp"}
	root.AddCommand(NewInitCmd())
	root.AddCommand(NewAddCmd())
	root.AddCommand(NewListCmd())
	root.AddCommand(NewGetCmd())
	root.AddCommand(NewRemoveCmd())
	root.AddCommand(NewEditCmd())
	root.AddCommand(NewExportCmd())
	root.AddCommand(NewImportCmd())
	root.AddCommand(NewPasswdCmd())

	return root
}

func TestCLIIntegration(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "gotp-cli-integ-*")
	defer os.RemoveAll(tmpDir)
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	// 1. Init
	t.Log("Testing Init")
	root := setupTestCLI(vaultPath, "password\npassword\n")
	_, err := executeCommand(root, "init")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 2. Add
	t.Log("Testing Add")
	root = setupTestCLI(vaultPath, "password\nJBSWY3DPEHPK3PXP\n\n\n")
	_, err = executeCommand(root, "add", "TestAccount", "--tags", "work")
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// 3. Edit
	t.Log("Testing Edit")
	root = setupTestCLI(vaultPath, "password\n")
	_, err = executeCommand(root, "edit", "TestAccount", "--name", "NewName")
	if err != nil {
		t.Fatalf("Edit failed: %v", err)
	}

	// 4. List advanced
	t.Log("Testing List Advanced")
	root = setupTestCLI(vaultPath, "password\n")
	out, _ := executeCommand(root, "list", "--filter", "work", "--sort", "issuer", "--with-codes")
	if !strings.Contains(out, "NewName") {
		t.Errorf("List advanced missing account. Got: %q", out)
	}

	// 9. Test Add Interactive
	t.Log("Testing Add Interactive")
	root = setupTestCLI(vaultPath, "password\nInteractiveAcc\nJBSWY3DPEHPK3PXP\n\n\n")
	_, err = executeCommand(root, "add")
	if err != nil {
		t.Fatalf("Interactive Add failed: %v", err)
	}

	// 10. Test Edit Interactive
	t.Log("Testing Edit Interactive")
	root = setupTestCLI(vaultPath, "password\n1\nEditedInteractive\n0\n")
	_, err = executeCommand(root, "edit", "InteractiveAcc")
	if err != nil {
		t.Fatalf("Interactive Edit failed: %v", err)
	}

	// 11. Test Get JSON
	t.Log("Testing Get JSON")
	root = setupTestCLI(vaultPath, "password\n")
	out, _ = executeCommand(root, "get", "EditedInteractive", "--json")
	if !strings.Contains(out, "code") {
		t.Errorf("Get JSON missing code. Got: %q", out)
	}

	// 12. Test Password Mismatch
	t.Log("Testing Password Mismatch")
	root = setupTestCLI(vaultPath, "password\nwrong\nwrong2\n")
	_, err = executeCommand(root, "passwd")
	if err == nil {
		t.Error("Expected error for password mismatch")
	}
}
