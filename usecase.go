package msgtm

type TagList interface {
	List() (*[]GitTag, error)
}

func MinorUpAll(tagList TagList) (*[]*ServiceTagWithSemVer, error) {
	tags, err := tagList.List()
	if err != nil {
		return nil, err
	}
	if len(*tags) == 0 {
		return nil, nil
	}

	serviceTags := []*ServiceTagWithSemVer{}
	for _, tag := range *tags {
		serviceTag, err := tag.ToServiceTag()
		if err != nil {
			continue
		}
		serviceTag.UpdateMinor()
		serviceTags = append(serviceTags, serviceTag)
	}

	return &serviceTags, nil
}

