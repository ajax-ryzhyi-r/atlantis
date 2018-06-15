package runtime

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-version"
	"github.com/runatlantis/atlantis/server/events/models"
)

// ApplyStepOperator runs `terraform apply`.
type ApplyStepOperator struct {
	TerraformExecutor TerraformExec
}

func (a *ApplyStepOperator) Run(ctx models.ProjectCommandContext, extraArgs []string, path string) (string, error) {
	planPath := filepath.Join(path, ctx.Workspace+".tfplan")
	stat, err := os.Stat(planPath)
	if err != nil || stat.IsDir() {
		return "", fmt.Errorf("no plan found at path %q and workspace %q–did you run plan?", ctx.RepoRelPath, ctx.Workspace)
	}

	tfApplyCmd := append(append(append([]string{"apply", "-no-color"}, extraArgs...), ctx.CommentArgs...), planPath)
	var tfVersion *version.Version
	if ctx.ProjectConfig != nil && ctx.ProjectConfig.TerraformVersion != nil {
		tfVersion = ctx.ProjectConfig.TerraformVersion
	}
	return a.TerraformExecutor.RunCommandWithVersion(ctx.Log, path, tfApplyCmd, tfVersion, ctx.Workspace)
}
