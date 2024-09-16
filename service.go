package msgtm

type TagList interface {
	List() (*[]GitTag, error)
}

func MajorUpAll(tagList TagList) (*[]*ServiceTagWithSemVer, error) {
	return versionUpAll(func(tag *ServiceTagWithSemVer) {
		tag.UpdateMajor()
	})(tagList)
}

func MinorUpAll(tagList TagList) (*[]*ServiceTagWithSemVer, error) {
	return versionUpAll(func(tag *ServiceTagWithSemVer) {
		tag.UpdateMinor()
	})(tagList)
}

func PatchUpAll(tagList TagList) (*[]*ServiceTagWithSemVer, error) {
	return versionUpAll(func(tag *ServiceTagWithSemVer) {
		tag.UpdatePatch()
	})(tagList)
}

type versionUpFunc func(*ServiceTagWithSemVer)

func versionUpAll(f versionUpFunc) VersionUpService {
	return func(tagList TagList) (*[]*ServiceTagWithSemVer, error) {
		tags, err := tagList.List()
		if err != nil {
			return nil, err
		}

		serviceTags := []*ServiceTagWithSemVer{}
		if tags == nil || len(*tags) == 0 {
			return &serviceTags, nil
		}

		tmpAlreadyUpdatedServiceTags := map[string]*ServiceTagWithSemVer{}

		for _, tag := range *tags {
			serviceTag, err := tag.ToServiceTag()
			if err != nil {
				continue
			}
			f(serviceTag)
			if version, ok := tmpAlreadyUpdatedServiceTags[serviceTag.service]; ok {
				if serviceTag.LessThan(version) || serviceTag.Equal(version) {
					continue
				}
				if serviceTag.GreaterThan(version) {
					tmpAlreadyUpdatedServiceTags[serviceTag.service] = serviceTag
					continue
				}
			}
			tmpAlreadyUpdatedServiceTags[serviceTag.service] = serviceTag
		}

		for _, serviceTag := range tmpAlreadyUpdatedServiceTags {
			serviceTags = append(serviceTags, serviceTag)
		}

		return &serviceTags, nil
	}
}
