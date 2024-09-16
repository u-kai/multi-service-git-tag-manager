package msgtm

type RegisterServiceTags interface {
	Register(*CommitId, *[]*ServiceTagWithSemVer) error
}

type TagList interface {
	List() (*[]GitTag, error)
}

type VersionUpServiceTag func(*[]GitTag) *[]*ServiceTagWithSemVer

func CreateServiceTags(
	registerService RegisterServiceTags,
	commitId *CommitId,
	serviceNames []string,
	version SemVer,
) error {
	serviceTags := []*ServiceTagWithSemVer{}
	for _, serviceName := range serviceNames {
		serviceTags = append(serviceTags, NewServiceTagWithSemVer(serviceName, version))
	}
	return registerService.Register(commitId, &serviceTags)
}

func VersionUpAllServiceTags(
	list TagList,
	registerService RegisterServiceTags,
	versionUpService VersionUpServiceTag,
	commitId *CommitId,
) error {
	tags, err := list.List()
	if err != nil {
		return err
	}

	updates := versionUpService(tags)

	err = registerService.Register(commitId, updates)
	if err != nil {
		return err
	}

	return nil
}
