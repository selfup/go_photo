package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	src := flag.String("src", "", "source directory (e.g., SD card/CFExpress Card/SSD/Directory on your Machine/etc..)")

	dst := flag.String("dst", "", "destination directory (e.g., ~/Pictures/Import on Mac, C:\\Users\\You\\Pictures\\Import on Windows)")

	wipe := flag.Bool("wipe", false, "delete source files after successful copy")

	tui := flag.Bool("tui", false, "launch interactive TUI mode")

	flag.Parse()

	if *tui {
		if err := runTUI(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)

			os.Exit(1)
		}

		return
	}

	if *src == "" || *dst == "" {
		fmt.Println("Usage: go_photo -src <source_dir> -dst <destination_dir> [-wipe]")
		fmt.Println("       go_photo -tui")

		os.Exit(1)
	}

	if err := run(*src, *dst, *wipe); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		os.Exit(1)
	}

	action := "copied"

	if *wipe {
		action = "moved"
	}

	fmt.Printf("All done! Your photos from %s have been %s to %s\n", *src, action, *dst)
}
