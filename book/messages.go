package book

type Borrow struct {
	ClientId uint32
	Id       uint32
}

type Return struct {
	ClientId uint32
	Id       uint32
}

type GetInformation struct {
}

type Information struct {
	response Book
}

type UnknownBook struct {
}

type NotAvailable struct {
}

type NewBook struct {
	Book Book
}
