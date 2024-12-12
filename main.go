package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	inputFile := flag.String("input", "", "path to yaml file to split")
	fileUrl := flag.String("url", "", "url to download yaml from")
	outDir := flag.String("out", ".", "output dir")

	flag.Parse()

	// Create out dir unless it's the local dir
	if *outDir != "." {
		err := os.MkdirAll(*outDir, 0755)
		if err != nil {
			fmt.Println("Error creating output directory:", err)
			os.Exit(1)
		}
	}

	// check if we have a local input or url to get the yaml from
	input := *inputFile
	var err error
	if *inputFile == "" {
		if *fileUrl == "" {
			fmt.Println("You must specify input or url")
			os.Exit(1)
		}
		input, err = getYaml(*fileUrl, *outDir)
		if err != nil {
			fmt.Println("couldn't download yaml", err)
			os.Exit(1)
		}
	}

	// read concatenated yaml
	file, err := os.Open(input)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	var inDocument bool

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			// Process the previous document
			if buffer.Len() > 0 {
				processDocument(buffer.String(), *outDir)
				buffer.Reset()
			}
			inDocument = true
		}
		if inDocument {
			buffer.WriteString(line + "\n")
		}
	}

	// Process the last document
	if buffer.Len() > 0 {
		processDocument(buffer.String(), *outDir)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	fmt.Println("files saved in", *outDir)
	os.Exit(0)
}

func processDocument(document, outDir string) {

	name := extractMetadataName(document)
	kind := extractKind(document)
	if name != "" {
		outputPath := filepath.Join(outDir, kind+"_"+name+".yaml")
		err := os.WriteFile(outputPath, []byte(document), 0644)
		if err != nil {
			fmt.Println("Error writing file:", err)
		}
	}

}

func extractMetadataName(document string) string {
	scanner := bufio.NewScanner(strings.NewReader(document))
	var foundMetadata bool
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "metadata:") {
			foundMetadata = true
		}
		if foundMetadata && strings.HasPrefix(strings.TrimSpace(line), "name:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				return parts[1]
			}
		}
	}
	return ""
}

func extractKind(document string) string {
	scanner := bufio.NewScanner(strings.NewReader(document))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "kind:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				return parts[1]
			}
		}
	}
	return ""
}

// getYaml takes url and a path in and downloads a yaml from the url, then saves it to the path
func getYaml(fileUrl string, outDir string) (string, error) {

	out, err := os.Create(outDir + "downloaded.yaml")
	if err != nil {
		return out.Name(), err
	}

	defer out.Close()

	resp, err := http.Get(fileUrl)
	if err != nil {
		return out.Name(), err
	}
	if resp.StatusCode != http.StatusOK {
		return out.Name(), fmt.Errorf("couldn't get file, received code %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return out.Name(), err
	}
	return outDir + out.Name(), nil

}
