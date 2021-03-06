package cmd

import (
	"fmt"

	"github.com/webhippie/redirects/model"
	"github.com/webhippie/redirects/store"
	"gopkg.in/urfave/cli.v2"
)

// Create provides the sub-command to create redirect patterns.
func Create() *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "Create a redirect pattern",
		ArgsUsage: " ",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "source",
				Value: "",
				Usage: "Source for the redirect",
			},
			&cli.StringFlag{
				Name:  "destination",
				Value: "",
				Usage: "Destination for the redirect",
			},
			&cli.IntFlag{
				Name:  "priority",
				Value: 0,
				Usage: "Priority for the redirect",
			},
		},
		Action: func(c *cli.Context) error {
			return Handle(c, handleCreate)
		},
	}
}

func handleCreate(c *cli.Context, s store.Store) error {
	record := &model.Redirect{}

	if val := c.String("source"); c.IsSet("source") && val != "" {
		record.Source = val
	} else {
		return fmt.Errorf("You must provide a source")
	}

	if val := c.String("destination"); c.IsSet("destination") && val != "" {
		record.Destination = val
	} else {
		return fmt.Errorf("You must provide a destination")
	}

	if val := c.Int("priority"); c.IsSet("priority") {
		record.Priority = val
	}

	return s.CreateRedirect(record)
}
