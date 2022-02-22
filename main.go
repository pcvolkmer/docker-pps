package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Process within docker container
type Process struct {
	ContainerID string
	Image       string
	PID         string
	UID         string
	Command     string
}

var (
	q    *bool
	h    *string
	uids *string
	help *bool
)

func init() {
	q = flag.BoolP("quiet", "q", false, "Only display process IDs")
	h = flag.StringP("host", "H", "", "Container `host`")
	uids = flag.String("uid", "", "Only display processes for `user`name/UID(s)")
	help = flag.Bool("help", false, "Show this help text")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "\nUsage: docker-pps [OPTIONS]\n\nShow list of Processes running in docker containers\n\nOptions:")
		flag.PrintDefaults()
	}

	flag.CommandLine.MarkHidden("help")
	flag.CommandLine.SortFlags = false

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	var opts []client.Opt
	if *h != "" {
		if !strings.Contains(*h, "//") {
			*h = "tcp://" + *h
		}
		if u, err := client.ParseHostURL(*h); err == nil {
			if u.Port() != "" {
				opts = append(opts, client.WithHost(*h))
			} else {
				opts = append(opts, client.WithHost(*h+":2375"))
			}
		} else {
			log.Fatal(err)
		}
	}

	dc, err := client.NewClientWithOpts(opts...)
	if err != nil {
		log.Fatalln(err)
		return
	}

	var processes []Process

	if containers, err := dc.ContainerList(context.Background(), types.ContainerListOptions{All: true}); err == nil {
		for _, container := range containers {
			tops, err := dc.ContainerTop(context.Background(), container.ID, []string{})
			if err != nil {
				continue
			}

			uid := 0
			pid := 0
			cmd := 0
			for idx, title := range tops.Titles {
				if title == "UID" {
					uid = idx
				}
				if title == "PID" {
					pid = idx
				}
				if title == "CMD" {
					cmd = idx
				}
			}

			for _, process := range tops.Processes {
				processes = append(processes, Process{
					ContainerID: container.ID,
					Image:       container.Image,
					PID:         process[pid],
					UID:         process[uid],
					Command:     process[cmd],
				})
			}
		}
	} else {
		log.Fatalln(err)
	}

	if *q == true {
		printQuietTable(processes)
		return
	}
	printTable(processes)
}

// Print out table if processes, running within docker containers
func printTable(processes []Process) {
	imageColumnLength := 5
	for _, process := range processes {
		if len(process.Image) > imageColumnLength {
			imageColumnLength = len(process.Image)
		}
	}

	fmt.Printf("CONTAINER ID   IMAGE%s   PID        UID        COMMAND\n", strings.Repeat(" ", imageColumnLength-5))

	for _, process := range processes {
		if len(*uids) == 0 || contains(strings.Split(*uids, ","), process.UID) {
			fmt.Printf("%8s   %s   %-8s   %-8s   %s\n",
				process.ContainerID[:12],
				process.Image+strings.Repeat(" ", imageColumnLength-len(process.Image)),
				process.PID,
				process.UID,
				process.Command)
		}
	}
}

// Print out list of IDs for processes, running within docker containers
func printQuietTable(processes []Process) {
	for _, process := range processes {
		fmt.Println(process.PID)
	}
}

func contains(uids []string, uid string) bool {
	for _, _uid := range uids {
		if _uid == uid {
			return true
		}
	}
	return false
}
