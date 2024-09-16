package msgtm

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
	if err != nil {
		return err
	}

	err = registerService.Register(commitId, updates)
	if err != nil {
		return err
	}

	return nil
}

type RegisterServiceTags interface {
	Register(*CommitId, *[]*ServiceTagWithSemVer) error
}

type TagList interface {
	List() (*[]GitTag, error)
}

type VersionUpServiceTag func(*[]GitTag) *[]*ServiceTagWithSemVer
