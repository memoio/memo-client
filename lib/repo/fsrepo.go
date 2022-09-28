package repo

import (
	"log"
	"os"
	"path/filepath"

	"github.com/memoio/memo-client/lib/backend/keystore"
	"github.com/memoio/memo-client/lib/types"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/xerrors"
)

const (
	keyStorePathPrefix = "keystore"
)

type FSRepo struct {
	path string

	keyDs types.KeyStore
	// lockfile io.Closer
}

func Exists(repoPath string) (bool, error) {
	_, err := os.Stat(filepath.Join(repoPath, keyStorePathPrefix))
	notExist := os.IsNotExist(err)
	if notExist {
		err = nil
	}
	return !notExist, err
}

func NewFSRepo(dir string) (*FSRepo, error) {
	repoPath, err := homedir.Expand(dir)
	if err != nil {
		return nil, err
	}

	if repoPath == "" {
		repoPath = "./"
	}

	err = ensureWritableDirectory(repoPath)
	if err != nil {
		return nil, xerrors.Errorf("no writable directory %w", err)
	}

	info, err := os.Stat(repoPath)
	if err != nil {
		return nil, xerrors.Errorf("failed to stat repo %s %w", repoPath, err)
	}

	var actualPath string
	if info.IsDir() {
		actualPath = repoPath
	} else {
		actualPath, err = os.Readlink(repoPath)
		if err != nil {
			return nil, xerrors.Errorf("failed to follow repo symlink %s %w", repoPath, err)
		}
	}

	r := &FSRepo{
		path: actualPath,
	}

	err = r.loadFromDisk()
	if err != nil {
		return nil, err
	}

	log.Println("open repo at:", repoPath)

	return r, nil
}

func ensureWritableDirectory(path string) error {
	// Attempt to create the requested directory, accepting that something might already be there.
	err := os.Mkdir(path, 0775)

	if err == nil {
		return nil // Skip the checks below, we just created it.
	} else if !os.IsExist(err) {
		return xerrors.Errorf("failed to create directory %s %w", path, err)
	}

	// Inspect existing directory.
	stat, err := os.Stat(path)
	if err != nil {
		return xerrors.Errorf("failed to stat path %s %w", path, err)
	}
	if !stat.IsDir() {
		return xerrors.Errorf("%s is not a directory", path)
	}
	if (stat.Mode() & 0600) != 0600 {
		return xerrors.Errorf("insufficient permissions for path %s, got %04o need %04o", path, stat.Mode(), 0600)
	}
	return nil
}

func (r *FSRepo) loadFromDisk() error {
	err := r.openKeyStore()
	if err != nil {
		return xerrors.Errorf("failed to open keystore %w", err)
	}
	return nil
}

func (r *FSRepo) openKeyStore() error {
	ksp := filepath.Join(r.path, "keystore")

	ks, err := keystore.NewKeyRepo(ksp)
	if err != nil {
		return err
	}

	r.keyDs = ks

	return nil
}

func (r *FSRepo) Close() error {
	err := r.keyDs.Close()
	if err != nil {
		return xerrors.Errorf("failed to close key store %w", err)
	}
	return nil
}

func (r *FSRepo) KeyStore() types.KeyStore {
	return r.keyDs
}

func (r *FSRepo) Path() (string, error) {
	return r.path, nil
}

func (r *FSRepo) Repo() Repo {
	return r
}
