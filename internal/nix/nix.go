package nix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ProfileElement struct {
	Active      bool     `json:"active"`
	AttrPath    string   `json:"attrPath"`
	OriginalURL string   `json:"originalUrl"`
	Outputs     any      `json:"outputs"`
	Priority    int      `json:"priority"`
	StorePaths  []string `json:"storePaths"`
	URL         string   `json:"url"`
}

type ProfileList struct {
	Elements map[string]ProfileElement `json:"elements"`
	Version  int                        `json:"version"`
}

func runNix(args ...string) (string, error) {
	cmd := exec.Command("nix", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("nix %s: %w\n%s", strings.Join(args, " "), err, stderr.String())
	}
	return stdout.String(), nil
}

func runNixPassthrough(args ...string) error {
	cmd := exec.Command("nix", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("nix %s: %w", strings.Join(args, " "), err)
	}
	return nil
}

func ProfileListJSON(profile string) (*ProfileList, error) {
	args := []string{"profile", "list", "--profile", profile, "--json"}
	out, err := runNix(args...)
	if err != nil {
		return nil, err
	}
	var list ProfileList
	if err := json.Unmarshal([]byte(out), &list); err != nil {
		return nil, fmt.Errorf("parsing nix profile list output: %w", err)
	}
	return &list, nil
}

func ProfileAdd(profile, ref string, priority int) error {
	return runNixPassthrough("profile", "add", "--profile", profile, ref, "--priority", fmt.Sprintf("%d", priority))
}

func ProfileRemove(profile, name string) error {
	return runNixPassthrough("profile", "remove", "--profile", profile, name)
}

func ProfileUpgrade(profile string, names []string, refresh bool) error {
	args := []string{"profile", "upgrade", "--profile", profile}
	if refresh {
		args = append(args, "--refresh")
	}
	args = append(args, names...)
	return runNixPassthrough(args...)
}

type LockNode struct {
	Inputs any `json:"inputs,omitempty"`
	Locked *struct {
		LastModified int64  `json:"lastModified"`
		NarHash     string `json:"narHash"`
		Owner       string `json:"owner,omitempty"`
		Repo        string `json:"repo,omitempty"`
		Rev         string `json:"rev"`
		Type        string `json:"type"`
		URL         string `json:"url,omitempty"`
	} `json:"locked,omitempty"`
	Original *struct {
		Owner string `json:"owner,omitempty"`
		Repo  string `json:"repo,omitempty"`
		Type  string `json:"type,omitempty"`
		URL   string `json:"url,omitempty"`
		Ref   string `json:"ref,omitempty"`
	} `json:"original,omitempty"`
}

type Locks struct {
	Root  string               `json:"root"`
	Nodes map[string]LockNode `json:"nodes"`
}

type FlakeMetadata struct {
	Fingerprint  string `json:"fingerprint"`
	LastModified int64  `json:"lastModified"`
	Locks        Locks  `json:"locks"`
}

func FlakeMetadataJSON(flake string) (*FlakeMetadata, error) {
	args := []string{"flake", "metadata", flake, "--json"}
	out, err := runNix(args...)
	if err != nil {
		return nil, err
	}
	var meta FlakeMetadata
	if err := json.Unmarshal([]byte(out), &meta); err != nil {
		return nil, fmt.Errorf("parsing nix flake metadata output: %w", err)
	}
	return &meta, nil
}

func FlakeUpdate(flake string, inputs []string, refresh bool) error {
	args := []string{"flake", "update", "--flake", flake}
	if refresh {
		args = append(args, "--refresh")
	}
	args = append(args, inputs...)
	return runNixPassthrough(args...)
}