package utils

import (
	"bytes"
	"fmt"
	"github.com/cmd-tools/gtfocli/constants"
	"github.com/cmd-tools/gtfocli/logger"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
	"path/filepath"
)

var GTFOBinCategories = []string{
	"file-read",
	"file-write",
	"shell",
	"sudo",
	"file-upload",
	"file-download",
	"limited-suid",
	"library-load",
	"capabilities",
}

type Category struct {
	Code        string `yaml:"code"        json:"code"        xml:"code"       `
	Description string `yaml:"description" json:"description" xml:"description"`
}

type Function struct {
	Name         string     `yaml:"name"          json:"name"           xml:"name"`
	FileRead     []Category `yaml:"file-read"     json:"file-read"      xml:"file-read"`
	FileWrite    []Category `yaml:"file-write"    json:"file-write"     xml:"file-write"`
	Shell        []Category `yaml:"shell"         json:"shell"          xml:"shell"`
	Sudo         []Category `yaml:"sudo"          json:"sudo"           xml:"sudo"`
	FileUpload   []Category `yaml:"file-upload"   json:"file-upload"    xml:"file-upload"`
	FileDownload []Category `yaml:"file-download" json:"file-download"  xml:"file-download"`
	LimitedSuid  []Category `yaml:"limited-suid"  json:"limited-suid"   xml:"limited-suid"`
	LibraryLoad  []Category `yaml:"library-load"  json:"library-load"   xml:"library-load"`
	Capabilities []Category `yaml:"capabilities"  json:"capabilities"   xml:"capabilities"`
}

type GTFOBinContent struct {
	Functions Function `yaml:"functions" json:"functions" xml:"functions"`
}

func CountGTFOBinsFiles() int {
	directoryPath := constants.GTFOBinsOutputDir
	files, err := os.ReadDir(directoryPath)
	logger.Logger.Infof("Get GTFOBin directory: %s", directoryPath)
	if err != nil {
		return 0
	}

	counter := 0
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != fmt.Sprintf(".%s", constants.GTFOBinsExtensions) {
			continue
		}
		counter++
	}
	logger.Logger.Infof("Found %d files in directory: %s", counter, directoryPath)

	return counter
}

func parseGTFOBinsFile(filePath string) (*Function, error) {
	yamlContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file %s: %v", filePath, err)
	}

	var gtfoBinContent GTFOBinContent
	err = yaml.Unmarshal(yamlContent, &gtfoBinContent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML from file %s: %v", filePath, err)
	}

	return &gtfoBinContent.Functions, nil
}

func GetGTFOBinsList() (map[string]*Function, error) {
	directoryPath := constants.GTFOBinsOutputDir
	files, err := os.ReadDir(directoryPath)
	logger.Logger.Infof("Get GTFOBin directory: %s", directoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %v", directoryPath, err)
	}

	functionsMap := make(map[string]*Function)
	logger.Logger.Infof("Found: %d files.", len(files))

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != fmt.Sprintf(".%s", constants.GTFOBinsExtensions) {
			continue
		}

		filePath := filepath.Join(directoryPath, file.Name())
		function, err := parseGTFOBinsFile(filePath)
		if err != nil {
			log.Printf("Error parsing YAML file %s: %v", filePath, err)
			continue
		}

		base := path.Base(file.Name())

		function.Name = base[0 : len(base)-len(path.Ext(base))]

		functionsMap[file.Name()] = function
	}

	return functionsMap, nil
}

func GetGTFOBinCategory(category string, list []*Function) []Category {
	var selectedCategories []Category
	for _, w := range list {
		switch category {
		case "shell":
			selectedCategories = append(selectedCategories, w.Shell...)
		case "sudo":
			selectedCategories = append(selectedCategories, w.Sudo...)
		case "file-read":
			selectedCategories = append(selectedCategories, w.FileRead...)
		case "file-write":
			selectedCategories = append(selectedCategories, w.FileWrite...)
		case "file-upload":
			selectedCategories = append(selectedCategories, w.FileUpload...)
		case "file-download":
			selectedCategories = append(selectedCategories, w.FileDownload...)
		case "limited-suid":
			selectedCategories = append(selectedCategories, w.LimitedSuid...)
		case "library-load":
			selectedCategories = append(selectedCategories, w.LibraryLoad...)
		case "capabilities":
			selectedCategories = append(selectedCategories, w.Capabilities...)
		}
	}
	return selectedCategories
}

func SummaryGTFOBinTable(data []*Function) bytes.Buffer {
	var buffer bytes.Buffer

	if len(data) == 0 {
		return buffer
	}

	table := tablewriter.NewWriter(&buffer)
	table.SetHeader(append([]string{"Name"}, GTFOBinCategories...))
	for _, item := range data {
		row := []string{
			item.Name,
			fmt.Sprintf("%d", len(item.FileRead)),
			fmt.Sprintf("%d", len(item.FileWrite)),
			fmt.Sprintf("%d", len(item.Shell)),
			fmt.Sprintf("%d", len(item.Sudo)),
			fmt.Sprintf("%d", len(item.FileUpload)),
			fmt.Sprintf("%d", len(item.FileDownload)),
			fmt.Sprintf("%d", len(item.LimitedSuid)),
			fmt.Sprintf("%d", len(item.LibraryLoad)),
			fmt.Sprintf("%d", len(item.Capabilities)),
		}
		table.Append(row)
	}
	table.Render()

	return buffer
}
