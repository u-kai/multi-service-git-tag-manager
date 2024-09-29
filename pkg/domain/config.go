package domain

import "fmt"

type ConfigFileService interface {
	Read() (*ConfigFile, error)
	Write(*ConfigFile) error
	Copy() *ConfigFile
}

type ConfigFile struct {
	Services *[]ServiceConfig
}

type ServiceConfig struct {
	Name   string
	Latest string
	Prev   string
	Desc   string
}

func (c *ConfigFile) UpdateServiceLatest(name *string, tag *GitTag, commitId *CommitId) {
	for _, service := range *c.Services {
		if service.Name == *name {
			service.Latest = tag.String()
			return
		}
	}
}

func (c *ConfigFile) UpdateServicePrev(name *string, tag *GitTag, commitId *CommitId) {
	for _, service := range *c.Services {
		if service.Name == *name {
			service.Prev = tag.String()
			return
		}
	}
}

type CommitInfo struct {
	CommitId CommitId
	Tag      GitTag
}

type ServiceTagList interface {
	List(string) ([]*CommitInfo, error)
}

func UpdateConfigFile(
	configFile ConfigFileService,
	serviceTagList ServiceTagList,
) error {
	config, err := configFile.Read()
	if err != nil {
		return fmt.Errorf("failed to read config file: %s", err.Error())
	}
	if config == nil {
		return fmt.Errorf("no services in config")
	}
	newConfigFile := configFile.Copy()
	for _, service := range *config.Services {
		tags, err := serviceTagList.List(service.Name)
		if err != nil {
			return fmt.Errorf("failed to get tags: %s", err.Error())
		}
		if len(tags) == 1 {
			newConfigFile.UpdateServiceLatest(&service.Name, &tags[0].Tag, &tags[0].CommitId)
			continue
		}
		if len(tags) > 2 {
			latestIndex := len(tags) - 1
			prevIndex := len(tags) - 2
			newConfigFile.UpdateServiceLatest(&service.Name, &tags[latestIndex].Tag, &tags[latestIndex].CommitId)
			newConfigFile.UpdateServicePrev(&service.Name, &tags[prevIndex].Tag, &tags[prevIndex].CommitId)
		}
	}
	err = configFile.Write(newConfigFile)
	if err != nil {
		return fmt.Errorf("failed to write config file: %s", err.Error())
	}
	return nil
}
