package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.i3wm.org/i3"
)

type Config struct {
	Separator   string            `json:"Separator"`
	DefaultIcon string            `json:"DefaultIcon"`
	Icons       map[string]string `json:"Icons"`
}

var config *Config

func main() {
	loadConfig()

	updateWorkspaceNames()

	wndEvent := i3.Subscribe(i3.WindowEventType)
	for wndEvent.Next() {
		ev := wndEvent.Event().(*i3.WindowEvent)
		switch ev.Change {
		case "new":
			fallthrough
		case "close":
			fallthrough
		case "move":
			updateWorkspaceNames()
		}
	}
	err := wndEvent.Close()
	if err != nil {
		panic(err)
	}
}

func loadConfig() {
	fmt.Println("loading config from ~/.i3-workspace-iconizer")
	config = &Config{
		Separator:   " ",
		DefaultIcon: "ﬓ",
		Icons:       map[string]string{},
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("failed to get user-home-dir, using default config: " + err.Error())
		return
	}

	b, err := os.ReadFile(home + "/.i3-workspace-iconizer")
	if err != nil {
		fmt.Println("failed to read config, using default config: " + err.Error())
		return
	}

	err = json.Unmarshal(b, config)
	if err != nil {
		fmt.Println("failed to unmarshal config, using default config: " + err.Error())
		return
	}
}

func updateWorkspaceNames() {

	tree, err := i3.GetTree()
	if err != nil {
		panic(err)
	}

	for _, output := range tree.Root.Nodes {
		for _, position := range output.Nodes {
			if position.Name != "content" {
				continue
			}

			for _, workspace := range position.Nodes {
				if workspace.Name == "__i3_scratch" {
					continue
				}

				name := strings.Split(workspace.Name, " ")[0]
				wnds := GetWindows(workspace)
				for _, w := range wnds {
					icon, ok := config.Icons[w.WindowProperties.Class]
					if ok {
						name += config.Separator + icon
					} else {
						name += config.Separator + config.DefaultIcon
					}
				}

				cmd := fmt.Sprintf("rename workspace \"%s\" to \"%s\"", workspace.Name, name)
				fmt.Println(cmd)
				r, err := i3.RunCommand(cmd)
				if err != nil {
					panic(err)
				}
				fmt.Println(r)
			}
		}
	}
}

func GetWindows(node *i3.Node) []*i3.Node {
	if len(node.Nodes) == 0 {
		return []*i3.Node{node}
	}
	result := []*i3.Node{}
	for _, n := range node.Nodes {
		result = append(result, GetWindows(n)...)
	}
	return result
}
