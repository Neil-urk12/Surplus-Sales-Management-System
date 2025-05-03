package models

type User struct {
	Id        string
	Name      string
	Email     string
	Password  string
	Role      string
	CreatedAt string
	UpdatedAt string
	IsActive  bool
	Token     string
}

type Customer struct {
	ID             string
	Name           string
	Email          string
	Phone          string
	Address        string
	DateRegistered string
	CreatedAt      string
	UpdatedAt      string
}

type Sale struct {
	ID         string
	CustomerID string
	SoldBy     string
	SaleDate   string
	TotalPrice float64
	CreatedAt  string
	UpdatedAt  string
}

type SaleItem struct {
	ID          string
	SaleID      string
	ItemType    string
	MultiCabID  string
	AccessoryID string
	MaterialID  string
	Quantity    int
	UnitPrice   float64
	Subtotal    float64
	CreatedAt   string
	UpdatedAt   string
}

type StockTransaction struct {
	ID             string
	UserID         string
	Timestamp      string
	Type           string
	QuantityChange int
	Remarks        string
	AccessoryID    string
	MaterialID     string
	CreatedAt      string
	UpdatedAt      string
}

type Accessory struct {
	ID        string
	Name      string
	Quantity  int
	CreatedAt string
	UpdatedAt string
}

type MultiCabAccessory struct {
	ID            string
	MultiCabID    string
	AccessoryID   string
	QuantityAdded int
	DateApplied   string
	CreatedAt     string
	UpdatedAt     string
}

type Material struct {
	ID        string
	Name      string
	Quantity  int
	CreatedAt string
	UpdatedAt string
}

type MultiCabMaterial struct {
	ID           string
	MultiCabID   string
	MaterialID   string
	QuantityUsed int
	DateApplied  string
	CreatedAt    string
	UpdatedAt    string
}

type MultiCab struct {
	ID           string
	Make         string
	Model        string
	Year         int
	Color        string
	Condition    string
	Price        float64
	Status       string
	DateAdded    string
	SerialNumber string
	CreatedAt    string
	UpdatedAt    string
}
