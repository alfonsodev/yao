package main

import (
	"fmt"

	_ "github.com/alfonsodev/yao/adapter/postgresql"
	g "github.com/alfonsodev/yao/generate"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func main() {
	type Flags struct {
		verbose   bool
		host      string
		database  string
		user      string
		pass      string
		sslmode   string
		outputDir string
	}

	var flags Flags

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Display version",
		Long:  `Display version of this software`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Yao (yet another orm) -- v0.3")
		},
	}

	var generateCmd = &cobra.Command{
		Use:   "gen",
		Short: "Generate models",
		Long:  `Generates one model file per each table in your datase.`,
		Run: func(cmd *cobra.Command, args []string) {
			params := fmt.Sprintf("dbname=%s sslmode=%s", flags.database, flags.sslmode)
			_, err := g.Open("postgres", params)

			if err != nil {
				fmt.Println(err.Error())
			}
			g.Generate("usermanager", flags.outputDir)
		},
	}

	generateCmd.Flags().StringVarP(&flags.outputDir, "output", "o", "", "Directory path where models folders will be created.")
	generateCmd.Flags().StringVarP(&flags.database, "database", "d", "", "Database name.")
	generateCmd.Flags().StringVarP(&flags.host, "host", "H", "", "Host name")
	generateCmd.Flags().StringVarP(&flags.user, "user", "u", "", "User name.")
	generateCmd.Flags().StringVarP(&flags.pass, "password", "p", "", "User password.")
	generateCmd.Flags().StringVarP(&flags.sslmode, "sslmode", "s", "disable", "ssl mode.")

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(versionCmd, generateCmd)
	rootCmd.Execute()
}
