package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/danielgatis/go-ruby-prism/parser"
)

const (
	railsRepoZipURL = "https://github.com/rails/rails/archive/refs/heads/main.zip"
	outputDir       = "test/fixtures"
	railsFolderName = "rails-main"
)

func main() {
	ctx := context.Background()

	// Step 1: Check if Rails source code is already downloaded
	zipPath := filepath.Join(outputDir, "rails.zip")
	unzippedPath := filepath.Join(outputDir, railsFolderName)

	if fileExists(zipPath) {
		fmt.Println("Rails source code zip already exists. Skipping download.")
	} else {
		fmt.Println("Downloading Rails source code...")
		if err := downloadFile(ctx, railsRepoZipURL, zipPath); err != nil {
			fmt.Println("Error downloading Rails source code:", err)
			return
		}
		fmt.Println("Rails source code downloaded to", zipPath)
	}

	// Step 2: Unzip the Rails folder if not already unzipped
	if dirExists(unzippedPath) {
		fmt.Println("Rails source code already unzipped. Skipping extraction.")
	} else {
		fmt.Println("Unzipping Rails source code...")
		if err := unzip(zipPath, outputDir); err != nil {
			fmt.Println("Error unzipping Rails source code:", err)
			return
		}
		fmt.Println("Rails source code unzipped to", unzippedPath)

		// Remove the zip file after unzipping
		if err := os.Remove(zipPath); err != nil {
			fmt.Println("Error removing zip file:", err)
			return
		}
		fmt.Println("Zip file removed:", zipPath)
	}

	// Step 3: Iterate over all Ruby files in the Rails folder
	if err := iterateRubyFiles(unzippedPath); err != nil {
		fmt.Println("Error iterating Ruby files:", err)
	}
}

// downloadFile downloads a file from the given URL and saves it to the specified path.
func downloadFile(ctx context.Context, url, dest string) error {
	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Download the file
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Write the response body to the file
	_, err = io.Copy(out, resp.Body)
	return err
}

// unzip extracts a zip file to the specified destination directory.
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		fpath := filepath.Join(dest, file.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// iterateRubyFiles iterates over all `.rb` files in the specified directory and prints their paths.
func iterateRubyFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		ctx := context.Background()
		p, _ := parser.NewParser(ctx)
		defer p.Close(ctx)

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".rb") {
			source, err := readFileContent(path)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			_, err = p.Parse(ctx, string(source))
			if err != nil {
				fmt.Println(err)
				fmt.Println("Ruby file with error:", path)
				os.Exit(1)
			}

			fmt.Println("Ruby file found:", path)
			if err == nil {
				fmt.Println("Ruby file parsed:", path)
			}
		}
		return nil
	})
}

// fileExists checks if a file exists and is not a directory.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// dirExists checks if a directory exists.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// readFileContent reads the content of a file and returns it as a byte slice.
func readFileContent(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}
