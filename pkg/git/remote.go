package git

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// RemoteHTTPURL returns the likely URL for where the repository in question
// has a web interface. It only supports HTTP and the SCP-like variant of the
// SSH protocol for the remote.
func (r *Repository) RemoteHTTPURL() (string, error) {
	url, err := url.Parse(r.Config.Remote.URL)
	// If the remote is just using HTTPS, then just return that minus the
	// `.git` suffix.
	if err == nil {
		if url.Scheme == "https" {
			return strings.TrimSuffix(r.Config.Remote.URL, ".git"), nil
		} else {
			return "", fmt.Errorf("unrecognized URL scheme in remote: %s", url.Scheme)
		}
	}

	// If that didn't work, let's assume it's using the SSH protocol.
	after, found := strings.CutPrefix(r.Config.Remote.URL, "git@")
	if !found {
		return "", errors.New("unsupported protocol for remote")
	}
	before, after, found := strings.Cut(after, ":")
	if !found {
		return "", errors.New("unsupported protocol for remote")
	}
	return fmt.Sprintf("https://%s/%s", before, strings.TrimSuffix(after, ".git")), nil
}
