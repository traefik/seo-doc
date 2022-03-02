package sitemap

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ldez/go-git-cmd-wrapper/v2/add"
	"github.com/ldez/go-git-cmd-wrapper/v2/commit"
	"github.com/ldez/go-git-cmd-wrapper/v2/config"
	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/ldez/go-git-cmd-wrapper/v2/push"
	"github.com/ldez/go-git-cmd-wrapper/v2/status"
	"github.com/urfave/cli/v2"
)

const defaultBranch = "master"

const (
	fileNameSitemap   = "sitemap.xml"
	fileGZNameSitemap = fileNameSitemap + ".gz"
)

// GitInfo represents the Git user configuration used for commit.
type GitInfo struct {
	UserName  string
	UserEmail string
	Token     string
}

// NewGitInfo creates a new GitInfo.
func NewGitInfo(cliCtx *cli.Context) GitInfo {
	return GitInfo{
		UserName:  cliCtx.String(flagGitUserName),
		UserEmail: cliCtx.String(flagGitUserEmail),
		Token:     cliCtx.String(flagGithubToken),
	}
}

// Commit commits and push the changes.
func Commit(gCfg GitInfo, debug bool) error {
	ctx := context.Background()

	// setup git user info
	output, err := setupGitUserInfo(gCfg, debug)
	if err != nil {
		fmt.Println(output)
		return fmt.Errorf("failed to set Git user: %w", err)
	}

	// check the git status of the dir
	output, err = git.StatusWithContext(ctx, status.Porcelain(""), git.Debugger(debug))
	if err != nil {
		fmt.Println(output)
		return fmt.Errorf("failed to get Git status: %w", err)
	}

	if !hasDiff(output) {
		log.Println("Nothing to commit.")
		return nil
	}

	// add target doc path to the index
	output, err = git.AddWithContext(ctx, add.PathSpec(fileNameSitemap, fileGZNameSitemap), git.Debugger(debug))
	if err != nil {
		log.Println(output)
		return fmt.Errorf("failed to add files: %w", err)
	}

	// create a commit
	output, err = git.CommitWithContext(ctx, commit.Message("Update sitemap files"), git.Debugger(debug))
	if err != nil {
		log.Println(output)
		return fmt.Errorf("failed to commit: %w", err)
	}

	// push the branch to the target git repo
	output, err = git.PushWithContext(ctx,
		push.Remote("origin"), push.RefSpec(defaultBranch),
		push.Repo(fmt.Sprintf("https://%s:@github.com/traefik/doc.git", gCfg.Token)),
		git.Debugger(debug))
	if err != nil {
		log.Println(output)
		return fmt.Errorf("failed to push: %w", err)
	}

	return nil
}

func setupGitUserInfo(gitInfo GitInfo, debug bool) (string, error) {
	if len(gitInfo.UserEmail) != 0 {
		output, err := git.Config(config.Entry("user.email", gitInfo.UserEmail), git.Debugger(debug))
		if err != nil {
			return output, err
		}
	}

	if len(gitInfo.UserName) != 0 {
		output, err := git.Config(config.Entry("user.name", gitInfo.UserName), git.Debugger(debug))
		if err != nil {
			return output, err
		}
	}

	return "", nil
}

func hasDiff(output string) bool {
	if len(output) == 0 {
		return false
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasSuffix(line, fileGZNameSitemap) || strings.HasSuffix(line, fileNameSitemap) {
			return true
		}
	}

	return false
}
