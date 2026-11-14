package upgrademng

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func createMockTarGz(t *testing.T, archivePath string, files map[string]string) {
	f, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("failed to create archive: %v", err)
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for name, content := range files {
		hdr := &tar.Header{
			Name: name,
			Mode: 0755,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatalf("failed to write tar header: %v", err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write tar content: %v", err)
		}
	}
}

func TestVerifySHA256(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.bin")
	content := []byte("hello world upgrade package")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	h := sha256.Sum256(content)
	expectedHash := hex.EncodeToString(h[:])

	if err := verifySHA256(filePath, expectedHash); err != nil {
		t.Fatalf("verifySHA256 failed with correct hash: %v", err)
	}

	if err := verifySHA256(filePath, "invalid_hash_12345"); err == nil {
		t.Fatal("verifySHA256 should fail with wrong hash")
	}
}

func TestExtractAndFindScript(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test-pkg.tar.gz")

	files := map[string]string{
		"mota/scripts/install.sh": "#!/bin/bash\necho install",
		"mota/bin/mota":          "binary content",
	}
	createMockTarGz(t, pkgPath, files)

	extractDir := filepath.Join(tmpDir, "extracted")
	if err := extractPackage(pkgPath, extractDir); err != nil {
		t.Fatalf("extractPackage failed: %v", err)
	}

	scriptPath, err := findInstallScript(extractDir, "mota")
	if err != nil {
		t.Fatalf("findInstallScript failed: %v", err)
	}

	expectedScript := filepath.Join(extractDir, "mota", "scripts", "install.sh")
	if scriptPath != expectedScript {
		t.Fatalf("expected script %s, got %s", expectedScript, scriptPath)
	}
}
