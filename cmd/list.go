package cmd

import (
    "io/ioutil"
    "fmt"
    "log"
    "github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command {
    Use: "list",
    Run: func(cmd *cobra.Command, args []string) {
        ListConnections()
    },
}

func ListConnections() {
    files, err := ioutil.ReadDir(configDir)
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
        name := file.Name()
        name = name[:len(name) - 4]
        fmt.Println(name)
    }
}
