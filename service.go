package msgtm

func MajorUpAll(tags *[]GitTag) *[]*ServiceTagWithSemVer {
	return versionUpAll(func(tag *ServiceTagWithSemVer) {
		tag.UpdateMajor()
	})(tags)
}

func MinorUpAll(tags *[]GitTag) *[]*ServiceTagWithSemVer {
	return versionUpAll(func(tag *ServiceTagWithSemVer) {
		tag.UpdateMinor()
	})(tags)
}

func PatchUpAll(tags *[]GitTag) *[]*ServiceTagWithSemVer {
	return versionUpAll(func(tag *ServiceTagWithSemVer) {
		tag.UpdatePatch()
	})(tags)
}

type versionUpFunc func(*ServiceTagWithSemVer)

func versionUpAll(f versionUpFunc) VersionUpServiceTag {
	return func(tags *[]GitTag) *[]*ServiceTagWithSemVer {
		serviceTags := []*ServiceTagWithSemVer{}
		if tags == nil || len(*tags) == 0 {
			return &serviceTags
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

		return &serviceTags
	}
}
