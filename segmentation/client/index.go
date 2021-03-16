package client

type Client interface {
	GetItems(offset int) (Response, error)
}
