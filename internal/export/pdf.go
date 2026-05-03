package export

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// ExportPDF converts paper.md to paper.pdf in finalDir using pandoc + wkhtmltopdf.
// It prints installation instructions if either tool is missing.
func ExportPDF(finalDir string) error {
	paperMD := filepath.Join(finalDir, "paper.md")
	paperPDF := filepath.Join(finalDir, "paper.pdf")

	if _, err := os.Stat(paperMD); os.IsNotExist(err) {
		return fmt.Errorf("paper.md not found at %s; run 'paperflow finalize' first", paperMD)
	}

	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		printInstallHint(paperMD)
		return fmt.Errorf("pandoc not found in PATH")
	}

	wk, err := exec.LookPath("wkhtmltopdf")
	if err != nil {
		printInstallHint(paperMD)
		return fmt.Errorf("wkhtmltopdf not found in PATH")
	}

	_ = wk

	cmd := exec.Command(pandoc, paperMD, "--pdf-engine=wkhtmltopdf", "-o", paperPDF)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pandoc: %w", err)
	}

	fmt.Printf("PDF written to: %s\n", paperPDF)
	return nil
}

func printInstallHint(paperMD string) {
	fmt.Fprintln(os.Stderr, "PDF export requires pandoc and wkhtmltopdf.")
	switch runtime.GOOS {
	case "windows":
		fmt.Fprintln(os.Stderr, "  winget install JohnMacFarlane.Pandoc")
		fmt.Fprintln(os.Stderr, "  winget install wkhtmltopdf.wkhtmltopdf")
	case "darwin":
		fmt.Fprintln(os.Stderr, "  brew install pandoc")
		fmt.Fprintln(os.Stderr, "  brew install wkhtmltopdf")
	default:
		fmt.Fprintln(os.Stderr, "  apt install pandoc wkhtmltopdf")
	}
	fmt.Fprintf(os.Stderr, "Your paper.md is ready at: %s\n", paperMD)
}
