package domain

type ItemStatus int

const (
	ItemStatusInitial ItemStatus = iota + 1
	ItemStatusOnSale
	ItemStatusSoldOut
)

type Item struct {
	ID          int32
	Name        string
	Price       int64
	Description string
	CategoryID  int64
	UserID      int64
	Image       []byte
	Status      ItemStatus
	CreatedAt   string
	UpdatedAt   string
}

type Category struct {
	ID   int64
	Name string
}
