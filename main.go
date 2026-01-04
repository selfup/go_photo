package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var extensionToFolder = map[string]string{
	".jpg":  "JPEG",
	".jpeg": "JPEG",
	".heif": "HEIF",
	".heic": "HEIF",
	".raw":  "RAW",
	".arw":  "RAW",
	".raf":  "RAW",
	".nef":  "RAW",
	".mov":  "MOV",
	".braw": "BRAW",
	".mp4":  "MP4",
}

func main() {
	src := flag.String("src", "", "source directory (e.g., SD card/CFExpress Card/SSD/Directory on your Machine/etc..)")

	dst := flag.String("dst", "", "destination directory (e.g., ~/Pictures/Import on Mac, C:\\Users\\You\\Pictures\\Import on Windows)")

	wipe := flag.Bool("wipe", false, "delete source files after successful copy")

	flag.Parse()

	if *src == "" || *dst == "" {
		fmt.Println("Usage: go_photo -src <source_dir> -dst <destination_dir> [-wipe]")

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

func run(src, dst string, wipe bool) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))

		folder, ok := extensionToFolder[ext]

		if !ok {
			return nil
		}

		destDir := filepath.Join(dst, folder)

		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", destDir, err)
		}

		destPath := filepath.Join(destDir, info.Name())

		if err := copyFile(path, destPath, info.ModTime()); err != nil {
			return fmt.Errorf("copying %s to %s: %w", path, destPath, err)
		}

		if wipe {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("deleting source file %s: %w", path, err)
			}

			fmt.Printf("Moved: %s -> %s\n", path, destPath)
		} else {
			fmt.Printf("Copied: %s -> %s\n", path, destPath)
		}

		return nil
	})
}

func copyFile(src, dst string, modTime time.Time) error {
	srcFile, err := os.Open(src)

	if err != nil {
		return err
	}

	defer srcFile.Close()

	dstFile, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	if err := dstFile.Sync(); err != nil {
		return err
	}

	return os.Chtimes(dst, modTime, modTime)
}
