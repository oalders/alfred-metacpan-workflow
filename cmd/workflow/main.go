package main

import (
	"fmt"
	"os"

	wf "github.com/oalders/alfred-metacpan-workflow"
	"github.com/urfave/cli"
)

var version = "dev"

func main() {
	app := cli.NewApp()
	app.Name = "Alfred Metacpan Workflow"
	app.Version = version
	app.Usage = ""
	app.Author = "handlename"
	app.Email = "nagata{at}handlena.me"
	app.Action = doMain
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "query",
			Value: "",
			Usage: "search query",
		},
	}
	app.Run(os.Args)
}

func doMain(c *cli.Context) {
	q := c.String("query")

	if len(q) < 3 {
		return
	}

	fmt.Println(wf.SearchModule(q))
}
