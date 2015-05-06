package main

import (
	"fmt"
	_ "github.com/alfonsodev/yao/adapter/postgresql"
	g "github.com/alfonsodev/yao/generate"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func main() {

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Display version",
		Long:  `Display version of this software`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Yao (yet another orm) -- v.01")
		},
	}

	var generateCmd = &cobra.Command{
		Use:   "gen",
		Short: "Generate models",
		Long:  `Generates one model file per each table in your datase.`,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := g.Open("postgres", "dbname=foo sslmode=disable")
			if err != nil {
				fmt.Println(err.Error())
			}
			//      g.Conn("postgres", db)
			g.Generate("usermanager")
		},
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(versionCmd, generateCmd)
	rootCmd.Execute()
}
