package providers

type Alpha struct {
  clientId, clientSecret string
}

func NewAlphaProvider(clientId, clientSecret string) Alpha {
  return Alpha{ clientId: clientId, clientSecret: clientSecret }
}

func (alpha Alpha) GetBalance() int64 {
  return 1312
}
