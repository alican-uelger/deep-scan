package e2e

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/alican-uelger/deep-scan/internal/scanner"
)

func createTestFiles(root string, numFiles int) error {
	for i := 0; i < numFiles; i++ {
		filePath := filepath.Join(root, "file"+strconv.Itoa(i)+".txt")
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			return err
		}
	}
	return nil
}

func BenchmarkSearch(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "search_test")
	if err != nil {
		b.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	numFiles := 1000 // Adjust the number of files for different load tests
	if err := createTestFiles(tempDir, numFiles); err != nil {
		b.Fatalf("failed to create test files: %v", err)
	}

	searcher := scanner.NewOs()
	options := scanner.SearchOptions{
		LogLate: false,
	}

	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := searcher.Search(tempDir, options)
		if err != nil {
			b.Fatalf("search failed: %v", err)
		}
	}
	b.StopTimer()
}
