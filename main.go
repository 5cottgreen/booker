package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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
					if c.Args().Len() < 2 {
						return errors.New("must have category name and task description")
					}

					os.Mkdir("tasks", os.ModeDir)
					tasksFile, err := os.OpenFile("tasks/"+c.Args().First()+".txt", os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						return err
					}

					defer tasksFile.Close()

					id, err := generateId()
					if err != nil {
						return err
					}

					tasksFile.Write([]byte(id + " " + c.Args().Slice()[1] + "\n"))
					fmt.Printf("task %s created\n", id)
					return nil
				},
			},

			{
				Name:  "done",
				Usage: "Romoves task.",
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.Args().Len() < 2 {
						return errors.New("must have category name and task id")
					}

					id := c.Args().Slice()[1]
					if id == "*" {
						if err := os.Remove("tasks/" + c.Args().First() + ".txt"); err != nil {
							fmt.Printf("task %s removed\n", id)
							return nil
						}

						fmt.Printf("tasks removed\n")
						return nil
					}

					tasksFile, err := os.Open("tasks/" + c.Args().First() + ".txt")
					if err != nil {
						fmt.Printf("task %s removed\n", id)
						return nil
					}

					content, err := io.ReadAll(tasksFile)
					if err != nil {
						return err
					}

					tasksFile.Close()
					if err := os.Remove("tasks/" + c.Args().First() + ".txt"); err != nil {
						return err
					}

					sep := []byte(" ")
					nsep := []byte("\n")
					tasks := bytes.Split(content, nsep)

					if len(tasks) == 2 {
						arg := string(bytes.SplitN(tasks[0], sep, 2)[0])
						if arg == id || arg == "" {
							fmt.Printf("task %s removed\n", id)
							return nil
						}
					}

					tasksFile, err = os.Create("tasks/" + c.Args().First() + ".txt")
					if err != nil {
						return err
					}

					defer tasksFile.Close()

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
					if c.Args().Len() > 0 {
						listTasks("tasks/" + c.Args().First() + ".txt")
						fmt.Println("linsting complete")
						return nil
					} else {
						items, err := os.ReadDir("tasks")
						if err != nil {
							fmt.Println("linsting complete")
							return nil
						}

						for _, item := range items {
							fmt.Println("|-" + strings.Split(item.Name(), ".")[0])
							listTasks("tasks/" + item.Name())
						}

						fmt.Println("linsting complete")
						return nil
					}
				},
			},

			{
				Name:  "article",
				Usage: "Creates article.",
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.Args().Len() < 1 {
						return errors.New("must have article name")
					}

					id, err := generateId()
					if err != nil {
						return err
					}

					os.Mkdir("articles", os.ModeDir)
					articleFile, err := os.OpenFile("articles/"+id+"_"+c.Args().First()+".txt", os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						return err
					}

					defer articleFile.Close()

					scanner := bufio.NewScanner(os.Stdin)
					for scanner.Scan() {
						articleFile.Write([]byte(scanner.Text() + "\n"))
					}

					if err := scanner.Err(); err != nil {
						return err
					}

					fmt.Println()
					fmt.Printf("article %s created\n", id)
					return nil
				},
			},

			{
				Name:  "delete",
				Usage: "Deletes article.",
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.Args().Len() < 1 {
						return errors.New("must have article id")
					}

					items, err := os.ReadDir("articles")
					if err != nil {
						fmt.Printf("article %s deleted\n", c.Args().First())
						return nil
					}

					for _, item := range items {
						if strings.SplitN(item.Name(), "_", 2)[0] == c.Args().First() {
							if err := os.Remove("articles/" + item.Name()); err != nil {
								fmt.Printf("article %s deleted\n", c.Args().First())
								return nil
							}

							fmt.Printf("article %s deleted\n", c.Args().First())
							return nil
						}
					}

					fmt.Printf("article %s deleted\n", c.Args().First())
					return nil
				},
			},

			{
				Name:  "guide",
				Usage: "Lists articles.",
				Action: func(ctx context.Context, c *cli.Command) error {
					items, err := os.ReadDir("articles")
					if err != nil {
						fmt.Println("linsting complete")
						return nil
					}

					if c.Args().Len() > 0 {
						for _, item := range items {
							if strings.SplitN(item.Name(), "_", 2)[0] == c.Args().First() {
								articleFile, err := os.Open("articles/" + item.Name())
								if err != nil {
									return err
								}

								content, err := io.ReadAll(articleFile)
								if err != nil {
									return err
								}

								fmt.Print(string(content))
								return nil
							}
						}
					} else {
						for _, item := range items {
							keys := strings.SplitN(item.Name(), "_", 2)
							fmt.Println(keys[0], strings.Split(keys[1], ".")[0])
						}

						fmt.Println()
						fmt.Println("linsting complete")
					}

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

// listTasks lists tasks
func listTasks(name string) (err error) {
	tasksFile, err := os.Open(name)
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

	return nil
}
