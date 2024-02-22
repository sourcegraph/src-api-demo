package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/sourcegraph/log"
	"github.com/sourcegraph/run"
	"github.com/sourcegraph/sourcegraph/lib/errors"
	"github.com/sourcegraph/src-api-demo/internal/observability"
	"github.com/sourcegraph/src-api-demo/internal/output"
	"github.com/urfave/cli/v2"
)

func main() {
	liblog := observability.InitLogs("gen", "dev")
	defer liblog.Sync()

	sort.Sort(cli.CommandsByName(gen.Commands))
	sort.Sort(cli.FlagsByName(gen.Flags))

	if err := gen.Run(os.Args); err != nil {
		_ = output.Render(output.FormatText, err)
		os.Exit(1)
	}
}

var gen = &cli.App{
	Name: "gen",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "version",
			Usage:    "The version of the Sourcegraph application to download the schema from",
			Required: false,
			Value:    "main",
		},
		// this is useful in bazel where all generator can share the same repo cache
		&cli.StringFlag{
			Name:  "archive.path",
			Usage: "The path to the repository archive to generate the schema from. version flag is ignored if this is set.",
		},
	},
	Action: func(c *cli.Context) error {
		logger := log.Scoped("gen").With(log.String("version", c.String("version")))

		cwd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "failed to get current working directory")
		}

		tmpDir, err := os.MkdirTemp("", "srcconf")
		if err != nil {
			return errors.Wrap(err, "failed to get user cache dir")
		}
		defer os.RemoveAll(tmpDir)

		version := c.String("version")
		repoDir := filepath.Join(tmpDir, fmt.Sprintf("sourcegraph-%s", version))
		defer os.RemoveAll(repoDir)

		if err := os.MkdirAll(repoDir, 0755); err != nil {
			return errors.Wrap(err, "failed to create repo dir")
		}

		archivePath := c.String("archive.path")
		if archivePath == "" {
			archivePath = filepath.Join(tmpDir, fmt.Sprintf("sourcegraph-%s.tar.gz", version))
			logger.Debug("downloading archive for version %s")
			if err := downloadFile(archivePath, fmt.Sprintf("https://github.com/sourcegraph/sourcegraph/archive/%s.tar.gz", version)); err != nil {
				return errors.Wrap(err, "failed to download archive")
			}
		}

		ctx := observability.LogCommands(c.Context, logger)
		// we use bsdtar instead of tar to avoid occasional "tar: dir: Directory renamed before its status could be extracted"
		if err := run.Cmd(ctx, "bsdtar", "-xzf", archivePath, "-C", repoDir, "--strip-components=1").Run().Wait(); err != nil {
			return errors.Wrap(err, "failed to extract archive")
		}

		schemaPath := filepath.Join(repoDir, "schema", "schema.go")
		schemaDstPath := filepath.Join(cwd, "schema.go")
		if err := run.Cmd(ctx, "cp", schemaPath, schemaDstPath).Run().Wait(); err != nil {
			return errors.Wrap(err, "failed to copy schema.go")
		}
		switch runtime.GOOS {
		case "darwin":
			err = run.Cmd(ctx, `sed -i '' -e 's/package schema/package srcconf/' ./schema.go`).Run().Wait()
		case "linux":
			err = run.Cmd(ctx, `sed -i 's/package schema/package srcconf/' ./schema.go`).Run().Wait()
		}
		if err != nil {
			return errors.Wrap(err, "failed to update schema.go")
		}

		return nil
	},
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Newf("failed to download file: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
