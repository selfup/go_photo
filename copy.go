package main

import (
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

type CopyProgress struct {
	Current   int
	Total     int
	File      string
	Done      bool
	Error     error
}

func countFiles(src string) (int, error) {
	count := 0

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))

		if _, ok := extensionToFolder[ext]; ok {
			count++
		}

		return nil
	})

	return count, err
}

func runWithProgress(src, dst string, wipe bool, progress chan<- CopyProgress) {
	defer close(progress)

	total, err := countFiles(src)

	if err != nil {
		progress <- CopyProgress{Error: err, Done: true}
		return
	}

	current := 0

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
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
		}

		current++

		progress <- CopyProgress{
			Current: current,
			Total:   total,
			File:    info.Name(),
			Done:    false,
		}

		return nil
	})

	if err != nil {
		progress <- CopyProgress{Error: err, Done: true}
		return
	}

	progress <- CopyProgress{Current: current, Total: total, Done: true}
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
