package msgtm

type CommitId string

func (c *CommitId) String() string {
	return string(*c)
}

const HEAD CommitId = "HEAD"

type RemoteAddr string

const Origin RemoteAddr = "origin"
