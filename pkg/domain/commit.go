package domain

type CommitId string

func (c *CommitId) String() string {
	return string(*c)
}

type RemoteAddr string

func (r *RemoteAddr) String() string {
	return string(*r)
}

const (
	HEAD   CommitId   = "HEAD"
	Origin RemoteAddr = "origin"
)
