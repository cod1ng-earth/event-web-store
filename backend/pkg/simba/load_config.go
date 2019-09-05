package simba

type contextDescription struct {
	Name         string
	Batch        bool
	ReadLock     string
	MessageNames []string
	Bridges      []subscriber
}

type subscriber struct {
	Name         string
	PkgPath      string
	MessageNames []string
}

type publisher struct {
	MessageNames []string
}
