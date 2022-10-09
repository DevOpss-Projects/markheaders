package markheaders

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// =======================================
// = Author : Mohamed Sherby (sn0walk3r) =
// = Twitter: @sn0walk3r                 =
// =======================================

func Entrypoint() {

	scan(".")
}
func readfile(file string) (string, error) {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(content), nil

}
func generateAttachmentHeaders(content string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))

	result := ""
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "![") {
			attachmentPath := strings.Split(strings.Split(scanner.Text(), "(")[1], ")")[0]

			result += fmt.Sprintf("<!-- Attachment: %s -->\n", attachmentPath)
		} else {
			continue
		}
	}
	//fmt.Println(result)
	return result, nil

}
func generateTitle(path string) string {
	paths := strings.Split(path, "\\")
	title := strings.Replace(paths[len(paths)-1], ".md", "", -1)
	return fmt.Sprintf("<!-- Title: %s -->\n", title)
}
func generateParentsHeaders(path string) ([]string, error) {
	parents := strings.Split(path, "\\")
	parents = parents[0 : len(parents)-1]
	for i := 0; i < len(parents); i++ {
		parents[i] = fmt.Sprintf("<!-- Parent: %s -->\n", parents[i])
	}
	return parents, nil
}
func modifyHeader(path string) (string, error) {
	result := ""
	content, file_err := readfile(path)
	if file_err != nil {
		return "", file_err
	}
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "<!--") {
			continue
		} else {
			result += scanner.Text() + "\n"
		}
	}
	fmt.Println(result)
	attachments, attachments_err := generateAttachmentHeaders(result)
	if attachments_err != nil {
		return "", attachments_err
	}

	parents, parents_err := generateParentsHeaders(path)
	if parents_err != nil {
		return "", parents_err
	}

	title := generateTitle(path)
	fmt.Printf("%s", title)

	content = strings.Join(parents, "") + attachments + title + result
	err := ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return "", err
	}

	return content, nil
}
func walkHandler(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if filepath.Ext(path) == ".md" {
		modifyHeader(path)
	}
	return nil
}
func scan(path string) {
	err := filepath.Walk(path, walkHandler)
	if err != nil {
		fmt.Println(err)
	}

}
