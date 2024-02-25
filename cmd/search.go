package cmd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/cmd-tools/gtfocli/constants"
	"github.com/cmd-tools/gtfocli/logger"
	"github.com/cmd-tools/gtfocli/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   constants.Search,
	Short: "Search for binary",
	Long:  "Search bypass for binary",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init(constants.Search, IsDebug())

		if len(args) > 0 {
			needleList = append(needleList, args[0])
		} else if inputFile != constants.Empty {
			needleList, _ = utils.ReadFromFile(inputFile)
		} else {
			needleList, _ = utils.ReadFromStdin()
		}

		if err := validateFlags(); nil != err {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if err := validateRepositories(); nil != err {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		needleList := cleanList(needleList)

		logger.Logger.Infof("Needle list contains: %d item/s.", len(needleList))
		logger.Logger.Infof("Selected operating system: %s.", operatingSystem)
		logger.Logger.Infof("Selected category: %s.", category)
		logger.Logger.Infof("Selected output: %s.", outputFormat)

		switch operatingSystem {
		case constants.Unix:
			results := searchInGTFOBin(needleList, category)
			fmt.Println(applyFormat(results, outputFormat))
		case constants.Windows:
			results := searchInLOLBAS(needleList, category)
			fmt.Println(applyFormat(results, outputFormat))
		}

		inputFile = constants.Empty
	},
}

var inputFile string
var outputFormat string
var operatingSystem string
var category string
var needleList []string
var allowedOS = []string{constants.Unix, constants.Windows}
var allowedOutputFormats = []string{constants.Text, constants.Json, constants.Yaml, constants.Xml}
var allowedWindowsCategories = append(utils.LOLBASCategories, constants.Summary)
var allowedUnixCategories = append(utils.GTFOBinCategories, constants.Summary)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVarP(&outputFormat, "output", "o", constants.Text, fmt.Sprintf("Define output type (%s).", strings.Join(allowedOutputFormats, "|")))
	searchCmd.Flags().StringVarP(&inputFile, "file", "f", constants.Empty, "File path which contains binary list names, one per line.")
	searchCmd.Flags().StringVarP(&operatingSystem, "os", constants.Empty, constants.Unix, fmt.Sprintf("Operating system binary files originate from (%s).", strings.Join(allowedOS, "|")))
	searchCmd.Flags().StringVarP(&category, "category", "c", constants.Summary, fmt.Sprintf(
		"Single category to select.\nFor %s: (%s).\nFor %s: (%s).",
		constants.Unix,
		strings.Join(allowedUnixCategories, "|"),
		constants.Windows,
		strings.Join(allowedWindowsCategories, "|")))
}

func cleanList(needleList []string) []string {
	var results []string
	for _, needle := range needleList {
		if needle != constants.Empty {
			base := filepath.Base(filepath.FromSlash(needle))
			results = append(results, base)
		}
	}
	return results
}

func validateFlags() error {

	IsAllowedOS := utils.IsStringInList(operatingSystem, allowedOS)
	if !IsAllowedOS {
		return fmt.Errorf("error: Allowed operating systems are: %s", strings.Join(allowedOS, "|"))
	}

	IsOutputFormatAllowed := utils.IsStringInList(outputFormat, allowedOutputFormats)
	if !IsOutputFormatAllowed {
		return fmt.Errorf("error: Allowed output formats are: %s", strings.Join(allowedOutputFormats, "|"))
	}

	switch operatingSystem {
	case constants.Unix:
		IsAllowedCategoryForSelectedOs := utils.IsStringInList(category, allowedUnixCategories)
		if !IsAllowedCategoryForSelectedOs {
			return fmt.Errorf("error: Allowed categories for operating system '%s' are: %s", constants.Unix, strings.Join(allowedUnixCategories, "|"))
		}
		break
	case constants.Windows:
		IsAllowedCategoryForSelectedOs := utils.IsStringInList(category, allowedWindowsCategories)
		if !IsAllowedCategoryForSelectedOs {
			return fmt.Errorf("error: Allowed categories for operating system '%s' are: %s", constants.Windows, strings.Join(allowedWindowsCategories, "|"))
		}
		break
	}
	return nil
}

func validateRepositories() error {
	if operatingSystem == constants.Unix && utils.CountGTFOBinsFiles() == 0 || operatingSystem == constants.Windows && utils.CountLOLBASFiles() == 0 {
		return fmt.Errorf("did you run first '%s %s' ", constants.Main, constants.Update)
	}
	return nil
}

func searchInGTFOBin(needleList []string, category string) interface{} {
	var result []*utils.Function

	list, err := utils.GetGTFOBinsList()
	if err != nil {
		return nil
	}

	for _, needle := range needleList {
		logger.Logger.Infof("Looking for needle: %s", needle)
		binary := list[fmt.Sprintf("%s.%s", needle, constants.GTFOBinsExtensions)]
		if binary != nil {
			result = append(result, binary)
			logger.Logger.Infof("Found 1 result for needle: %s", needle)
		}
	}

	if category != constants.Summary {
		return utils.GetGTFOBinCategory(category, result)
	}

	return result
}

func searchInLOLBAS(needleList []string, category string) interface{} {
	var result []*utils.LOLBASContent

	list, err := utils.GetLOLBASList()
	if err != nil {
		return nil
	}

	for _, needle := range needleList {
		logger.Logger.Infof("Looking for needle: %s", needle)
		binary := list[fmt.Sprintf("%s.%s", needle, constants.LOLBASExtensions)]
		if binary != nil {
			result = append(result, binary)
			logger.Logger.Infof("Found 1 result for needle: %s", needle)
		}
	}

	if category != constants.Summary {
		return utils.GetLOLBASCategory(category, result)
	}

	return result
}

func applyFormat(output interface{}, format string) string {
	switch format {
	case constants.Json:
		marshal, err := json.Marshal(output)
		if err != nil {
			return constants.Empty
		}
		return string(marshal)
	case constants.Yaml:
		marshal, err := yaml.Marshal(output)
		if err != nil {
			return constants.Empty
		}
		return string(marshal)
	case constants.Xml:
		marshal, err := xml.Marshal(output)
		if err != nil {
			return constants.Empty
		}
		return string(marshal)
	case constants.Text:
		if operatingSystem == constants.Unix {
			functions, ok := output.([]*utils.Function)
			if !ok {
				categories, _ := output.([]utils.Category)
				for _, i := range categories {
					if i.Description == constants.Empty {
						i.Description = "not available."
					}
					fmt.Printf("Code: %s\nDescription: %s\n", i.Code, i.Description)
					fmt.Println("------------------------")
				}
				break

			}
			table := utils.SummaryGTFOBinTable(functions)
			return table.String()
		}

		if operatingSystem == constants.Windows {
			functionsSlice, ok := output.([]*utils.LOLBASContent)
			if !ok {
				commands, _ := output.([]utils.Commands)
				for _, i := range commands {
					if i.Description == constants.Empty {
						i.Description = "not available."
					}
					fmt.Printf("Code: %s\n"+
						"Description: %s\n"+
						"MitreID: %s\n"+
						"OperatingSystem: %s\n"+
						"Privileges: %s\n"+
						"Usecase: %s\n",
						i.Command,
						i.Description,
						i.MitreID,
						i.OperatingSystem,
						i.Privileges,
						i.Usecase)
					fmt.Println("------------------------")
				}
				break
			}
			table := utils.SummaryLOLBASTable(functionsSlice)
			return table.String()
		}
	}

	return constants.Empty
}
