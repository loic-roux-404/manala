package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/loaders"
	"manala/syncer"
	"manala/validator"
)

// UpdateCmd represents the update command
func UpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update [dir]",
		Aliases: []string{"up"},
		Short:   "Update project",
		Long: `Update (manala update) will update project, based on
recipe and related variables defined in manala.yaml.

Example: manala update -> resulting in an update in a directory (default to the current directory)`,
		Run:  updateRun,
		Args: cobra.MaximumNArgs(1),
	}

	return cmd
}

func updateRun(cmd *cobra.Command, args []string) {
	// Loaders
	repoLoader := loaders.NewRepositoryLoader(viper.GetString("cache_dir"))
	recLoader := loaders.NewRecipeLoader()
	prjLoader := loaders.NewProjectLoader(repoLoader, recLoader, viper.GetString("repository"))

	// Project directory
	dir := viper.GetString("dir")
	if len(args) != 0 {
		// Get directory from first command arg
		dir = args[0]
	}

	// Load project
	prj, err := prjLoader.Load(dir)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Validate project
	if err := validator.ValidateProject(prj); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Project validated")

	// Sync project
	if err := syncer.SyncProject(prj); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Project synced")
}
