package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	flag.Parse()

	if *src == "" || *dst == "" {
		fmt.Println("Usage: go_photo -src <source_dir> -dst <destination_dir>")
		os.Exit(1)
	}

	if err := run(*src, *dst); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("All done! Your photos from %s have been copied to %s\n", *src, *dst)
}

func run(src, dst string) error {
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
		if err := copyFile(path, destPath); err != nil {
			return fmt.Errorf("copying %s to %s: %w", path, destPath, err)
		}

		fmt.Printf("Copied: %s -> %s\n", path, destPath)
		return nil
	})
}

func copyFile(src, dst string) error {
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

	return dstFile.Sync()
}
