package endpoint

type IEndpoint interface {
	IEndpointClient
}

type IEndpointClient interface {
	GetClientToken(string) (string, error)
}