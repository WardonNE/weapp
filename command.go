package weapp

import "github.com/spf13/cobra"

func migrateCmd(app *Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "run database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			rollback, err := cmd.Flags().GetBool("rollback")
			if err != nil {
				return err
			}
			if rollback {
				app.migrations.Rollback()
			} else {
				app.migrations.Commit()
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolP("rollback", "r", false, "rollback migrations")
	return cmd
}
