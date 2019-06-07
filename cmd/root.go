package cmd

import (
    "gopkg.in/yaml.v2"
    "log"
    "io"
    "path"
    "os"
    "io/ioutil"
    "github.com/spf13/cobra"
    "os/exec"
)

type Config struct {
    Host string   `yaml:"host"`
    User string   `yaml:"user"`
    Domain string `yaml:"domain"`
    Password string `yaml:"password"`
}

var configDir = "/home/max/.config/rdp-cli"

var rootCmd = &cobra.Command {
    Use: "rdp-cli",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        Run(args[0])
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}

func Run(name string) {
    path := path.Join(configDir, name + ".yml")
    fd, err := os.Open(path)
    if err != nil {
        log.Fatalf("error: %v", err)
    }
    defer fd.Close()
    data, err := ioutil.ReadAll(fd)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    config := Config{}

    err = yaml.Unmarshal([]byte(data), &config)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    args := []string{"-k", "de", "-K", "-r", "clipboard:CLIPBOARD", "-g", "1920x1040"}

    if config.User != "" {
        args = append(args, "-u", config.User)
    }

    if config.Domain != "" {
        args = append(args, "-d", config.Domain)
    }

    pass := ""
    if config.Password != "" {
        passArgs := []string{"lookup", config.Password, "password"}
        passCmd := exec.Command("secret-tool", passArgs...)
        output, err := passCmd.Output()
        if err != nil {
            log.Fatal(err)
        }
        pass = string(output)
        args = append(args, "-p", "-")
    }

    args = append(args, config.Host)

    cmd := exec.Command("rdesktop", args...)
    if pass != "" {
        stdin, err := cmd.StdinPipe()
        if err != nil {
            log.Fatal(err)
        }
        go func() {
            defer stdin.Close()
            io.WriteString(stdin, pass)
        }()
    }
    err = cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
}
