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
	"path/filepath"
	"regexp"
	"strings"
)

var LOLBASCategories = []string{
	"ads",
	"awl-bypass",
	"compile",
	"copy",
	"credentials",
	"decode",
	"download",
	"dump",
	"encode",
	"execute",
	"reconnaissance",
	"uac-bypass",
	"upload",
}

type Acknowledgement struct {
	Person string `yaml:"Person" json:"Person" xml:"Person"`
	Handle string `yaml:"Handle" json:"Handle" xml:"Handle"`
}

type FullPath struct {
	Path string `yaml:"Path" json:"Path" xml:"Path"`
}

type Resources struct {
	Path string `yaml:"Link" json:"Link" xml:"Link"`
}

type Detection struct {
	Sigma     string `yaml:"Sigma" json:"Sigma" xml:"Sigma"`
	Elastic   string `yaml:"Elastic" json:"Elastic" xml:"Elastic"`
	Splunk    string `yaml:"Splunk" json:"Splunk" xml:"Splunk"`
	BlockRule string `yaml:"BlockRule" json:"BlockRule" xml:"BlockRule"`
	IOC       string `yaml:"IOC" json:"IOC" xml:"IOC"`
}

type Commands struct {
	Command         string `yaml:"Command"         json:"Command"         xml:"Command"`
	Description     string `yaml:"Description"     json:"Description"     xml:"Description"`
	Usecase         string `yaml:"Usecase"         json:"Usecase"         xml:"Usecase"`
	Category        string `yaml:"Category"        json:"Category"        xml:"Category"`
	Privileges      string `yaml:"Privileges"      json:"Privileges"      xml:"Privileges"`
	MitreID         string `yaml:"MitreID"         json:"MitreID"         xml:"MitreID"`
	OperatingSystem string `yaml:"OperatingSystem" json:"OperatingSystem" xml:"OperatingSystem"`
}

type LOLBASContent struct {
	Name            string            `yaml:"Name"           json:"Name"             xml:"Name"`
	Author          string            `yaml:"Author"          json:"Author"          xml:"Author"`
	Description     string            `yaml:"Description"     json:"Description"     xml:"Description"`
	Created         string            `yaml:"Created"         json:"Created"         xml:"Created"`
	Commands        []Commands        `yaml:"Commands"        json:"Commands"        xml:"Commands"`
	FullPath        []FullPath        `yaml:"Full_Path"       json:"Full_Path"       xml:"Full_Path"`
	Detection       []Detection       `yaml:"Detection"       json:"Detection"       xml:"Detection"`
	Resources       []Resources       `yaml:"Resources"       json:"Resources"       xml:"Resources"`
	Acknowledgement []Acknowledgement `yaml:"Acknowledgement" json:"Acknowledgement" xml:"Acknowledgement"`
}

func CountLOLBASFiles() int {
	directoryPath := constants.LOLBASOutputDir
	files, err := os.ReadDir(directoryPath)
	logger.Logger.Infof("Get LOLBAS directory: %s", directoryPath)
	if err != nil {
		return 0
	}

	counter := 0
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != fmt.Sprintf(".%s", constants.LOLBASExtensions) {
			continue
		}
		counter++
	}
	logger.Logger.Infof("Found %d files in directory: %s", counter, directoryPath)

	return counter
}

func parseLOLBASFile(filePath string) (*LOLBASContent, error) {
	yamlContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file %s: %v", filePath, err)
	}

	var functions LOLBASContent
	err = yaml.Unmarshal(yamlContent, &functions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML from file %s: %v", filePath, err)
	}

	return &functions, nil
}

func GetLOLBASList() (map[string]*LOLBASContent, error) {
	directoryPath := constants.LOLBASOutputDir
	files, err := os.ReadDir(directoryPath)
	logger.Logger.Infof("Get LOLBAS directory: %s", directoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %v", directoryPath, err)
	}

	functionsMap := make(map[string]*LOLBASContent)
	logger.Logger.Infof("Found: %d files.", len(files))

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != fmt.Sprintf(".%s", constants.LOLBASExtensions) {
			continue
		}

		filePath := filepath.Join(directoryPath, file.Name())
		functions, err := parseLOLBASFile(filePath)
		if err != nil {
			log.Printf("Error parsing YAML file %s: %v", filePath, err)
			continue
		}

		functions.Name = file.Name()

		functionsMap[file.Name()] = functions
	}

	return functionsMap, nil
}

func GetLOLBASCategory(category string, list []*LOLBASContent) []Commands {
	var selectedCommands []Commands
	for _, b := range list {
		for _, command := range b.Commands {
			if getFormattedCategory(command.Category) == category {
				selectedCommands = append(selectedCommands, command)
			}
		}
	}
	return selectedCommands
}

func SummaryLOLBASTable(data []*LOLBASContent) bytes.Buffer {
	var buffer bytes.Buffer

	if len(data) == 0 {
		return buffer
	}

	table := tablewriter.NewWriter(&buffer)
	table.SetHeader(append([]string{"Name"}, LOLBASCategories...))
	for _, item := range data {
		var list []string
		for _, category := range LOLBASCategories {
			list = append(list, fmt.Sprintf("%d", countCategory(category, item.Commands)))
		}
		table.Append(append([]string{item.Name}, list...))
	}
	table.Render()

	return buffer
}

func countCategory(category string, commands []Commands) int {
	sum := 0
	for _, command := range commands {
		if getFormattedCategory(command.Category) == category {
			sum++
		}
	}
	return sum
}

func getFormattedCategory(category string) string {
	re := regexp.MustCompile(`\s+`)
	formattedCategory := strings.ToLower(re.ReplaceAllString(category, " "))
	return strings.ReplaceAll(formattedCategory, " ", "-")
}
