package ui

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"strings"

	files "github.com/tung-dnt/pagoda/public"
)

// staticVersions maps each embedded static file, keyed by its path relative to public/static, to a
// short hash of its content. The files are embedded in the binary, so this is computed once at
// startup and only changes when a file's content does, not on every deploy or restart.
var staticVersions = hashStatic(files.Static)

// hashStatic hashes every file under the "static" directory of fsys.
func hashStatic(fsys fs.FS) map[string]string {
	versions := make(map[string]string)
	err := fs.WalkDir(fsys, "static", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		f, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}
		versions[strings.TrimPrefix(path, "static/")] = hex.EncodeToString(h.Sum(nil)[:5])
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("failed to hash static files: %v", err))
	}
	return versions
}

// PublicFile generates a relative URL to a public file.
func PublicFile(filepath string) string {
	return fmt.Sprintf("/%s/%s", "files", filepath)
}

// StaticFile generates a relative URL to a static file. Files that exist get a cache-buster query
// parameter derived from their content, so long-lived cache headers are safe.
func StaticFile(filepath string) string {
	if v, ok := staticVersions[filepath]; ok {
		return fmt.Sprintf("/%s/%s?v=%s", "static", filepath, v)
	}
	return fmt.Sprintf("/%s/%s", "static", filepath)
}
