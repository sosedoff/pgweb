package aws

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sosedoff/pgweb/pkg/command"

	"github.com/go-ini/ini"
)

func awsConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".aws/config")
}

func awsCredentialsPath() string {
	return filepath.Join(os.Getenv("HOME"), ".aws/credentials")
}

func readProfile(opts *command.Options) {
	log.Println("[aws] trying to load configuration from ~/.aws directory")

	// Read the config file
	config, err := ini.Load(awsConfigPath())
	if err != nil {
		log.Println("[aws] cant read config file:", err)
		return
	}

	// Read credentials file
	credentials, err := ini.Load(awsCredentialsPath())
	if err != nil {
		log.Println("[aws] cant read credentials file:", err)
		return
	}

	// Find AWS profile
	profileName := opts.AWSProfile
	if profileName != "" && profileName != "default" {
		profileName = "profile " + profileName
	}
	section, err := config.GetSection(profileName)
	if err != nil {
		log.Println("[aws] cant find profile:", err)
		return
	}

	// Assign region if its not set
	if opts.AWSRegion == "" {
		if key := section.Key("region"); key != nil {
			opts.AWSRegion = key.Value()
		}
	}

	// Assign credentials
	if section := credentials.Section(opts.AWSProfile); section != nil {
		for _, k := range section.Keys() {
			switch k.Name() {
			case "aws_access_key_id":
				opts.AWSAccessKey = k.Value()
			case "aws_secret_access_key":
				opts.AWSSecretKey = k.Value()
			}
		}
	} else {
		log.Println("[aws] cant find credentials profile:", opts.AWSProfile)
		return
	}
}
