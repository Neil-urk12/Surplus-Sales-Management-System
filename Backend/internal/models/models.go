package models

import "time"

type User struct {
	Id        string    `json:"id"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	IsActive  bool      `json:"isActive"`
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
	AccessoryID    int
	MaterialID     string
	CreatedAt      string
	UpdatedAt      string
}

type Accessory struct {
	ID        int       `json:"id"`         // Unique identifier
	Name      string    `json:"name"`       // Name of the accessory
	Make      string    `json:"make"`       // Manufacturer/brand of the accessory
	Quantity  int       `json:"quantity"`   // Number of units available
	Price     float64   `json:"price"`      // Price in PHP
	Status    string    `json:"status"`     // Inventory status
	UnitColor string    `json:"unit_color"` // Color of the accessory
	Image     string    `json:"image"`      // URL or base64 string of the image
	CreatedAt time.Time `json:"createdAt"`  // Timestamp of creation
	UpdatedAt time.Time `json:"updatedAt"`  // Timestamp of last update
}

type MultiCabAccessory struct {
	ID            string
	MultiCabID    string
	AccessoryID   int // Changed from string to int to match the Accessory.ID type
	QuantityAdded int
	DateApplied   string
	CreatedAt     string
	UpdatedAt     string
}

type Material struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Supplier  string    `json:"supplier"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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

// MultiCab defines the structure for cab data, aligning with frontend needs.
type MultiCab struct {
	ID        int       `json:"id"`         // Unique identifier
	Name      string    `json:"name"`       // Name of the cab model (e.g., RX-7)
	Make      string    `json:"make"`       // Manufacturer (e.g., Mazda)
	Quantity  int       `json:"quantity"`   // Number of units available
	Price     float64   `json:"price"`      // Price in PHP
	Status    string    `json:"status"`     // Inventory status (e.g., In Stock, Low Stock)
	UnitColor string    `json:"unit_color"` // Color of the cab unit
	Image     string    `json:"image"`      // URL or base64 string of the image
	CreatedAt time.Time `json:"createdAt"`  // Timestamp of creation
	UpdatedAt time.Time `json:"updatedAt"`  // Timestamp of last update
}
