package simba

type contextDescription struct {
	Name         string
	Batch        bool
	ReadLock     string
	MessageNames []string
	Bridges      []subscriber
	Publisher    publisher
}

type subscriber struct {
	Name         string
	PkgPath      string
	MessageNames []string
}

type publisher struct {
	PkgPath      string
	MessageNames []string
}
