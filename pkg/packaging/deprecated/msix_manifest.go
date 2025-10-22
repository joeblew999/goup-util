package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const manifestTemplate = `
<?xml version="1.0" encoding="utf-8"?>
<Package
  xmlns="http://schemas.microsoft.com/appx/manifest/foundation/windows10"
  xmlns:uap="http://schemas.microsoft.com/appx/manifest/uap/windows10"
  xmlns:rescap="http://schemas.microsoft.com/appx/manifest/foundation/windows10/restrictedcapabilities"
  IgnorableNamespaces="uap rescap">

  <Identity
    Name="{{.name}}"
    Publisher="{{.publisher}}"
    Version="{{.version}}" />

  <Properties>
    <DisplayName>{{.displayName}}</DisplayName>
    <PublisherDisplayName>{{.publisherDisplayName}}</PublisherDisplayName>
    <Logo>assets/logo.png</Logo>
  </Properties>

  <Dependencies>
    <TargetDeviceFamily Name="Windows.Desktop" MinVersion="10.0.17763.0" MaxVersionTested="10.0.19041.0" />
  </Dependencies>

  <Resources>
    <Resource Language="en-us" />
  </Resources>

  <Applications>
    <Application Id="App"
      Executable="{{.executable}}"
      EntryPoint="Windows.FullTrustApplication">
      <uap:VisualElements
        DisplayName="{{.displayName}}"
        Description="{{.description}}"
        BackgroundColor="transparent"
        Square150x150Logo="assets/Square150x150Logo.png"
        Square44x44Logo="assets/Square44x44Logo.png">
      </uap:VisualElements>
    </Application>
  </Applications>

  <Capabilities>
    <rescap:Capability Name="runFullTrust" />
  </Capabilities>
</Package>
`

// msixManifestCmd represents the msix-manifest command
// DEPRECATED: Use 'bundle windows' command instead
var msixManifestCmd = &cobra.Command{
	Use:        "msix-manifest",
	Short:      "Generate an MSIX manifest file from a template (DEPRECATED)",
	Deprecated: "use 'bundle windows' instead for unified packaging workflow",
	Long: `⚠️  DEPRECATED: This command is replaced by 'bundle windows'

Migration:
  Old: goup-util msix-manifest --output AppxManifest.xml --name MyApp
  New: goup-util bundle windows examples/myapp

The new bundle command provides:
- Unified workflow (same as macOS/Android/iOS)
- Automatic manifest generation from templates
- Integrated MSIX creation
- Better error handling

See: docs/PACKAGING.md for complete guide

---

Generates an MSIX manifest file (AppxManifest.xml) for a Windows application.

You can provide the manifest data using flags or by providing a YAML data file.
Values from flags will override values from the data file.`,
	Run: func(cmd *cobra.Command, args []string) {
		outputPath, _ := cmd.Flags().GetString("output")
		dataPath, _ := cmd.Flags().GetString("data")
		templatePath, _ := cmd.Flags().GetString("template")

		if outputPath == "" {
			fmt.Println("Please provide an output path for the manifest file.")
			return
		}

		// Start with an empty data map
		data := make(map[string]interface{})

		// Load data from file if provided
		if dataPath != "" {
			yamlFile, err := ioutil.ReadFile(dataPath)
			if err != nil {
				fmt.Printf("Error reading data file: %v\n", err)
				return
			}
			err = yaml.Unmarshal(yamlFile, &data)
			if err != nil {
				fmt.Printf("Error unmarshalling data file: %v\n", err)
				return
			}
		}

		// Override with flags
		if cmd.Flags().Changed("name") {
			data["name"], _ = cmd.Flags().GetString("name")
		}
		if cmd.Flags().Changed("publisher") {
			data["publisher"], _ = cmd.Flags().GetString("publisher")
		}
		if cmd.Flags().Changed("version") {
			data["version"], _ = cmd.Flags().GetString("version")
		}
		if cmd.Flags().Changed("display-name") {
			data["displayName"], _ = cmd.Flags().GetString("display-name")
		}
		if cmd.Flags().Changed("publisher-display-name") {
			data["publisherDisplayName"], _ = cmd.Flags().GetString("publisher-display-name")
		}
		if cmd.Flags().Changed("executable") {
			data["executable"], _ = cmd.Flags().GetString("executable")
		}
		if cmd.Flags().Changed("description") {
			data["description"], _ = cmd.Flags().GetString("description")
		}

		// Get the template content
		var tmplContent string
		if templatePath != "" {
			contentBytes, err := ioutil.ReadFile(templatePath)
			if err != nil {
				fmt.Printf("Error reading template file: %v\n", err)
				return
			}
			tmplContent = string(contentBytes)
		} else {
			tmplContent = manifestTemplate
		}

		// Parse and execute the template
		tmpl, err := template.New("manifest").Parse(tmplContent)
		if err != nil {
			fmt.Printf("Error parsing template: %v\n", err)
			return
		}

		outputFile, err := os.Create(outputPath)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			return
		}
		defer outputFile.Close()

		err = tmpl.Execute(outputFile, data)
		if err != nil {
			fmt.Printf("Error executing template: %v\n", err)
			return
		}

		fmt.Printf("Successfully generated MSIX manifest at %s\n", outputPath)
	},
}

func init() {
	rootCmd.AddCommand(msixManifestCmd)
	msixManifestCmd.Flags().StringP("output", "o", "", "Output path for the AppxManifest.xml file (required)")
	msixManifestCmd.Flags().StringP("template", "t", "", "Path to a custom AppxManifest.xml template file")
	msixManifestCmd.Flags().StringP("data", "D", "", "Path to a YAML file containing manifest data")

	msixManifestCmd.Flags().StringP("name", "n", "", "The name of the application (e.g., MyCompany.MyApp)")
	msixManifestCmd.Flags().StringP("publisher", "p", "", "The publisher of the application (e.g., CN=MyCompany)")
	msixManifestCmd.Flags().StringP("version", "v", "", "The version of the application (e.g., 1.0.0.0)")
	msixManifestCmd.Flags().StringP("display-name", "d", "", "The display name of the application")
	msixManifestCmd.Flags().StringP("publisher-display-name", "b", "", "The publisher display name")
	msixManifestCmd.Flags().StringP("executable", "e", "", "The name of the executable file")
	msixManifestCmd.Flags().StringP("description", "s", "", "The description of the application")
}
