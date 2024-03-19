// A module for handling module release in daggerverse

package main

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func New(
	// A git repository where the release process will be applied
	gitRepo *Directory,
	// The module name to publish
	component string,

) *ModReleaser {
	return &ModReleaser{
		Component: component,
		Ctr: dag.
			Container().
			From("alpine:latest").
			WithExec([]string{"apk", "add", "--no-cache", "git"}).
			WithDirectory("/opt/repo/", gitRepo).
			WithWorkdir("/opt/repo/"),
	}
}

type ModReleaser struct {
	Tags []string
	Tag  string
	// +private
	Ctr *Container
	// +private
	Component string
}

func (m *ModReleaser) ListTags(ctx context.Context) (*ModReleaser, error) {
	versionRegexp, err := regexp.Compile(m.Component + "/v\\d+\\.\\d+\\.\\d+")
	if err != nil {
		return nil, err
	}

	output, err := m.Ctr.WithExec([]string{"git", "tag", "-l"}).Stdout(ctx)
	if err != nil {
		return nil, err
	}

	for _, tag := range strings.Split(output, "\n") {
		if versionRegexp.MatchString(tag) {
			m.Tags = append(m.Tags, tag)
		}
	}

	slices.Sort(m.Tags)

	return m, nil
}

func (m *ModReleaser) WithGitConfig(cfg *File) *ModReleaser {
	return m
}

func (m *ModReleaser) WithGitConfigEmail(email string) *ModReleaser {
	return m
}

func (m *ModReleaser) WithGitConfigName(name string) *ModReleaser {
	return m
}

func (m *ModReleaser) WithBranch(
	// Define the branch fro where to publish
	// +optional
	// +default="main"
	branch string,
) *ModReleaser {
	return m.WithContainer(m.Ctr.WithExec([]string{"git", "checkout", branch}))
}

func (m *ModReleaser) Major(
	// Define a custom message for the git tag otherwise it will be the default from the function
	// +optional
	msg string,
) (*ModReleaser, error) {
	return m.bumpVersion(true, false, false, msg)
}

func (m *ModReleaser) Minor(
	// Define a custom message for the git tag otherwise it will be the default from the function
	// +optional
	msg string,
) (*ModReleaser, error) {
	return m.bumpVersion(false, true, false, msg)
}

func (m *ModReleaser) Patch(
	// Define a custom message for the git tag otherwise it will be the default from the function
	// +optional
	msg string,
) (*ModReleaser, error) {
	return m.bumpVersion(false, false, true, msg)
}

func (m *ModReleaser) Publish() (*ModReleaser, error) {
	m.Ctr.WithExec([]string{"git", "push", "origin", m.Tag})

	return nil, nil
}

func (m *ModReleaser) WithContainer(ctr *Container) *ModReleaser {
	m.Ctr = ctr

	return m
}

func (m *ModReleaser) bumpVersion(major, minor, patch bool, customMsg string) (*ModReleaser, error) {
	var msg string
	var firstRelease bool
	prefixTag := m.Component + "/v"
	if len(m.Tags) == 0 {
		msg = "New component: " + m.Component
		m.Tag = prefixTag + "0.1.0"
		firstRelease = true
	}

	if !firstRelease {
		verMajor, verMinor, verPatch, err := m.parseSemver(strings.Split(strings.Split(m.Tags[len(m.Tags)-1], "/v")[1], "."))
		if err != nil {
			return nil, err
		}

		switch {
		case major:
			verMajor++
		case minor:
			verMinor++
		case patch:
			verPatch++
		default:
			return nil, fmt.Errorf("'major', 'minor', or 'patch' should be set to true")
		}

		m.Tag = prefixTag + fmt.Sprintf("%d.%d.%d", verMajor, verMinor, verPatch)
		msg = "New release " + m.Tag
	}

	if customMsg != "" {
		msg = customMsg
	}

	return m.WithContainer(m.Ctr.WithExec([]string{"git", "tag", "-a", m.Tag, "-m", msg})), nil
}

func (m ModReleaser) parseSemver(semverParts []string) (int, int, int, error) {
	major, err := strconv.Atoi(semverParts[0])
	if err != nil {
		return 0, 0, 0, err
	}
	minor, err := strconv.Atoi(semverParts[1])
	if err != nil {
		return 0, 0, 0, err
	}
	patch, err := strconv.Atoi(semverParts[2])
	if err != nil {
		return 0, 0, 0, err
	}

	return major, minor, patch, nil
}

func (m ModReleaser) Shell() *Terminal {
	return m.Ctr.WithDefaultTerminalCmd(nil).Terminal()
}
