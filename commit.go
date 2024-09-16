package msgtm

type CommitId string

func (c *CommitId) String() string {
	return string(*c)
}

const HEAD CommitId = "HEAD"

type RemoteAddr string

func (r *RemoteAddr) String() string {
	return string(*r)
}

const Origin RemoteAddr = "origin"
