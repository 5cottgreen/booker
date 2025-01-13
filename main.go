package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "greet",
				Usage: "Greets you.",
				Action: func(ctx context.Context, c *cli.Command) error {
					fmt.Println("Hello friend!")
					return nil
				},
			},

			{
				Name:  "cheerup",
				Usage: "Cheers you up.",
				Action: func(ctx context.Context, c *cli.Command) error {
					fmt.Println("Never stop fighting!")
					return nil
				},
			},

			{
				Name:  "task",
				Usage: "Creates task.",
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.Args().Len() < 1 {
						return errors.New("must have at least one symbol")
					}

					tasksFile, err := os.OpenFile("tasks.txt", os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						return err
					}

					defer tasksFile.Close()

					id, err := generateId()
					if err != nil {
						return err
					}

					tasksFile.Write([]byte(id + " " + c.Args().Slice()[0] + "\n"))
					fmt.Printf("task %s created\n", id)
					return nil
				},
			},

			{
				Name:  "done",
				Usage: "Romoves task.",
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.Args().Len() < 1 {
						return errors.New("must have at least one symbol")
					}

					id := c.Args().Slice()[0]
					if id == "*" {
						if err := os.Remove("tasks.txt"); err != nil {
							return err
						}

						fmt.Printf("tasks removed\n")
						return nil
					}

					tasksFile, err := os.Open("tasks.txt")
					if err != nil {
						return err
					}

					content, err := io.ReadAll(tasksFile)
					if err != nil {
						return err
					}

					tasksFile.Close()
					if err := os.Remove("tasks.txt"); err != nil {
						return err
					}

					tasksFile, err = os.Create("tasks.txt")
					if err != nil {
						return err
					}

					defer tasksFile.Close()

					sep := []byte(" ")
					nsep := []byte("\n")
					tasks := bytes.Split(content, nsep)

					for _, task := range tasks {
						arg := string(bytes.SplitN(task, sep, 2)[0])
						if arg == id || arg == "" {
							continue
						}

						if _, err := tasksFile.Write(task); err != nil {
							return err
						}

						if _, err := tasksFile.Write(nsep); err != nil {
							return err
						}
					}

					fmt.Printf("task %s removed\n", id)
					return nil
				},
			},

			{
				Name:  "list",
				Usage: "Lists tasks.",
				Action: func(ctx context.Context, c *cli.Command) error {
					tasksFile, err := os.Open("tasks.txt")
					if err != nil {
						return err
					}

					defer tasksFile.Close()

					content, err := io.ReadAll(tasksFile)
					if err != nil {
						return err
					}

					tasks := bytes.Split(content, []byte("\n"))
					for _, task := range tasks {
						fmt.Println(string(task))
					}

					fmt.Println("linsting complete")
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

// generateId generates ID
func generateId() (id string, err error) {
	b := make([]byte, 2)
	_, err = rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
