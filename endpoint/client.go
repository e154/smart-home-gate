package endpoint

type Client struct {
	*CommonEndpoint
}

func NewClient(common *CommonEndpoint) *Client {
	return &Client{CommonEndpoint: common}
}

func (c *Client) GetClientToken(clientId string) (token string, err error) {

	return
}
