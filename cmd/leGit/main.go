package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func calculateSHA1(typeOfObject string, lenOfContent string, content string) ([]byte, string) {
	hash := sha1.New()
	hash.Write([]byte(typeOfObject + " " + lenOfContent + "\x00" + content))
	hashBytes := hash.Sum(nil)
	hashString := fmt.Sprintf("%x", hashBytes)

	return hashBytes, hashString
}

func compressData(typeOfObject string, lenOfContent string, content string) []byte {
	var compressedData bytes.Buffer
	zlibWriter := zlib.NewWriter(&compressedData)
	if _, err := zlibWriter.Write([]byte(typeOfObject + " " + lenOfContent + "\x00" + content)); err != nil {
		fmt.Fprintf(os.Stderr, "Error compressing file: %s\n", err)
		os.Exit(1)
	}
	zlibWriter.Close()

	return compressedData.Bytes()
}

func writeObjectToFile(hashString string, compressedData []byte) {
	objDir := ".git/objects/" + hashString[:2]
	if err := os.MkdirAll(objDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
		os.Exit(1)
	}

	objFilePath := objDir + "/" + hashString[2:]
	if err := os.WriteFile(objFilePath, compressedData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		os.Exit(1)
	}

}

func createBlobObject(fileName string) (string, []byte) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}

	// Calculate SHA1 hash of the content
	hashBytes, hashString := calculateSHA1("blob", strconv.Itoa(len(content)), string(content))

	// Compress the content
	compressedData := compressData("blob", strconv.Itoa(len(content)), string(content))

	// Write compressed content to file
	writeObjectToFile(hashString, compressedData)

	return hashString, hashBytes
}

func createTreeObject(dirPath string) (string, []byte) {
	dir, err := os.Open(dirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening directory: %s\n", err)
		os.Exit(1)
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("Error reading the directory:", err)
		os.Exit(1)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	var content string
	var size int

	for _, file := range files {
		if file.Mode().IsDir() {
			if file.Name() != ".git" {
				_, compressedData := createTreeObject(dirPath + "/" + file.Name())
				content += fmt.Sprintf("40000 %s\x00%s", file.Name(), compressedData)
				size += len("40000") + len(file.Name()) + 22
			}
		} else {
			_, compressedData := createBlobObject(dirPath + "/" + file.Name())
			code := "100644"
			content += fmt.Sprintf("%s %s\x00%s", code, file.Name(), compressedData)
			size += len("100644") + len(file.Name()) + 22
		}
	}

	// Calculate SHA1 hash of the content
	hashBytes, hashString := calculateSHA1("tree", strconv.Itoa(size), content)

	// Compress the content
	compressedData := compressData("tree", strconv.Itoa(size), content)

	// Write compressed content to file
	writeObjectToFile(hashString, compressedData)

	return hashString, hashBytes
}

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}
	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}
		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}
		fmt.Println("Initialized git directory")

	case "cat-file":
		blobHash := os.Args[3]
		blobDir := blobHash[:2]
		file, err := os.Open(".git/objects/" + blobDir + "/" + blobHash[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()

		reader := bufio.NewReader(file)

		zlibReader, err := zlib.NewReader(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating zlib reader: %s\n", err)
			os.Exit(1)
		}
		defer zlibReader.Close()

		decompressedData, err := io.ReadAll(zlibReader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading object content: %s\n", err)
			os.Exit(1)
		}

		content := string(decompressedData)

		index := -1
		for i, b := range content {
			if b == '\x00' {
				index = i
				break
			}
		}
		if index == -1 {
			fmt.Fprintf(os.Stderr, "Error empty file\n")
			os.Exit(1)
		}

		actualContent := content[index+1:]
		fmt.Printf("%s", actualContent)

	case "hash-object":
		fileName := os.Args[3]
		hashString, _ := createBlobObject(fileName)
		fmt.Printf("%s", hashString)

	case "ls-tree":
		treeObjectHash := os.Args[len(os.Args)-1]
		treeDir := treeObjectHash[:2]
		file, err := os.Open(".git/objects/" + treeDir + "/" + treeObjectHash[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()

		reader := bufio.NewReader(file)

		zlibReader, err := zlib.NewReader(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating zlib reader: %s\n", err)
			os.Exit(1)
		}
		defer zlibReader.Close()

		decompressedData, err := io.ReadAll(zlibReader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading object content: %s\n", err)
			os.Exit(1)
		}

		content := string(decompressedData)

		index := -1
		for i, b := range content {
			if b == '\x00' {
				index = i
				break
			}
		}
		if index == -1 {
			fmt.Fprintf(os.Stderr, "Error empty file\n")
			os.Exit(1)
		}

		actualContent := content[index+1:]

		if os.Args[2] == "--name-only" {
			entries := strings.Split(actualContent, "\x00")

			for _, entry := range entries {
				value := strings.Split(entry, " ")
				if len(value) >= 2 {
					name := value[1]
					fmt.Printf("%s\n", name)
				}
			}
		} else {
			entries := strings.Split(actualContent, "\x00")

			for i, entry := range entries {
				value := strings.Split(entry, " ")
				if len(value) >= 2 && i != len(entries)-1 {
					var sha1Hash, mode, typeOfObject string
					if i == 0 {
						mode = value[0]
					} else {
						mode = value[0][20:]
					}
					if mode == "04000" {
						typeOfObject = "tree"
					} else {
						typeOfObject = "blob"
					}
					name := value[1]
					valueTemp := strings.Split(entries[i+1], " ")
					if len(valueTemp) >= 2 {
						sha1Hash = valueTemp[0][:20]
					}
					fmt.Printf("%s %s %s    %s\n", mode, typeOfObject, sha1Hash, name)
				}
			}
		}

	case "write-tree":
		hashString, _ := createTreeObject(".")
		fmt.Printf("%s", hashString)

	case "commit-tree":
		treeHash := os.Args[2]
		flag := os.Args[3]
		var preCommitHash string
		if flag == "-p" {
			preCommitHash = os.Args[4]
		}
		commitMessage := os.Args[len(os.Args)-1]
		currentTime := time.Now()
		currentZone := currentTime.Format("+0000")
		var content string
		content += fmt.Sprintf("tree %s\n", treeHash)
		if flag == "-p" {
			content += fmt.Sprintf("parent %s\n", preCommitHash)
		}
		content += fmt.Sprintf("author Mohit <mohit.pal.singh@outlook.com> %x %s\n", currentTime, currentZone)
		content += fmt.Sprintf("committer Mohit <mohit.pal.singh@outlook.com> %x %s\n", currentTime, currentZone)
		content += fmt.Sprintf("\n%s\n", commitMessage)

		// Calculate SHA1 hash of the content
		_, hashString := calculateSHA1("commit", strconv.Itoa(len(content)), content)

		// Compress the content
		compressedData := compressData("commit", strconv.Itoa(len(content)), content)

		// Write compressed content to file
		writeObjectToFile(hashString, compressedData)
		fmt.Println(hashString)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
