package qdrantkit

type Qdrant struct {
	Host           string
	CollectionName string
}

func New(host, cn string) *Qdrant {
	q := &Qdrant{
		Host:           host,
		CollectionName: cn,
	}
	q.init()
	return q
}
