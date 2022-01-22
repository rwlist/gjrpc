package astinfo

import (
	"bytes"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// A Dir describes a directory holding code by specifying
// the expected import path and the file system directory.
type Dir struct {
	importPath string // import path for that dir
	dir        string // file system directory
}

func findModuleRoots() ([]Dir, error) {
	var list []Dir

	// Check for use of modules by 'go env GOMOD',
	// which reports a go.mod file path if modules are enabled.
	stdout, _ := exec.Command("go", "env", "GOMOD").Output()
	gomod := string(bytes.TrimSpace(stdout))

	if gomod == "" || gomod == os.DevNull {
		return nil, errors.Errorf("no go modules detected")
	}

	cmd := exec.Command("go", "list", "-m", "-f={{.Path}}\t{{.Dir}}", "all")
	cmd.Stderr = os.Stderr
	out, _ := cmd.Output()
	for _, line := range strings.Split(string(out), "\n") {
		i := strings.Index(line, "\t")
		if i < 0 {
			continue
		}
		path, dir := line[:i], line[i+1:]
		if dir != "" {
			list = append(list, Dir{importPath: path, dir: dir})
		}
	}

	return list, nil
}

// findPackagePath finds package import path by absolute directory path.
func findPackagePath(pkgDir string) (string, error) {
	dirs, err := findModuleRoots()
	if err != nil {
		return "", err
	}

	for _, root := range dirs {
		if pkgDir == root.dir {
			return root.importPath, nil
		}
		if strings.HasPrefix(pkgDir, root.dir+string(filepath.Separator)) {
			suffix := filepath.ToSlash(pkgDir[len(root.dir)+1:])
			if root.importPath == "" {
				return suffix, nil
			}
			return root.importPath + "/" + suffix, nil
		}
	}

	return "", errors.Errorf("not found")
}
