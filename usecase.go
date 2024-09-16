package msgtm

func VersionUpAllServiceTags(
	list TagList,
	registerService RegisterServiceTags,
	versionUpService VersionUpService,
	commitId *CommitId,
) error {
	updates, err := versionUpService(list)
	if err != nil {
		return err
	}
	err = registerService.Register(commitId, updates)
	if err != nil {
		return err
	}
	//	if updates == nil || len(*updates) == 0 {
	//		return nil
	//	}
	//	for _, update := range *updates {
	//		gitTag := update.ToGitTag()
	//		// このままである処理は成功してある処理は成功しない可能性がある
	//		// RDBじゃないので、それはしょうがないとして、何が成功したのか何が失敗したのかはerrで表現しても良いかも
	//		err := repository.AddTag(&gitTag, commitId)
	//		if err != nil {
	//			return err
	//		}
	//	}

	return nil
}

type VersionUpService func(tagList TagList) (*[]*ServiceTagWithSemVer, error)
