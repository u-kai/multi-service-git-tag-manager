package usecase

import "msgtm/pkg/domain"

func PushAll(
	commitGetter CommitTagGetter,
	pusher CommitPusher,
	remote *domain.RemoteAddr,
	commitId *domain.CommitId,
) error {
	tags, err := commitGetter.Execute(GetCommitTagQuery{CommitId: commitId})
	if err != nil {
		return err
	}

	serviceTags := domain.FilterServiceTags(tags)

	err = pusher.Execute(CommitPushCommand{
		RemoteAddr: remote,
		Tags:       serviceTags,
	})
	if err != nil {
		return err
	}

	return nil
}
