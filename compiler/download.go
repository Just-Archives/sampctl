package compiler

import (
	"fmt"
	"os"
	"runtime"

	"github.com/pkg/errors"

	"github.com/Southclaws/sampctl/download"
	"github.com/Southclaws/sampctl/util"
)

// FromCache attempts to get a compiler package from the cache, `hit` represents success
func FromCache(cacheDir string, version Version, dir string) (hit bool, err error) {
	pkg, filename, err := GetCompilerPackageInfo(runtime.GOOS, version)
	if err != nil {
		return false, err
	}

	hit, err = download.FromCache(cacheDir, filename, dir, pkg.Method, pkg.Paths)
	if !hit {
		return false, nil
	}

	fmt.Printf("Using cached package %s\n", filename)

	return
}

// FromNet downloads a compiler package to the cache
func FromNet(cacheDir string, version Version, dir string) (err error) {
	fmt.Printf("Downloading compiler package %s\n", version)

	pkg, filename, err := GetCompilerPackageInfo(runtime.GOOS, version)
	if err != nil {
		return errors.Wrap(err, "package info mismatch")
	}

	if !util.Exists(dir) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return errors.Wrapf(err, "failed to create dir %s", dir)
		}
	}

	if !util.Exists(cacheDir) {
		err := os.MkdirAll(cacheDir, 0755)
		if err != nil {
			return errors.Wrapf(err, "failed to create cache %s", cacheDir)
		}
	}

	path, err := download.FromNet(pkg.URL, cacheDir, filename)
	if err != nil {
		return errors.Wrap(err, "failed to download package")
	}

	err = pkg.Method(path, dir, pkg.Paths)
	if err != nil {
		return errors.Wrapf(err, "failed to unzip package %s", path)
	}

	return
}
