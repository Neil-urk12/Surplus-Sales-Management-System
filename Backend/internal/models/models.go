package models

import (
	"time"
)

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	IsActive  bool      `json:"isActive"`
}

type Customer struct {
	ID             string    `json:"id"`
	FullName       string    `json:"fullName"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	DateRegistered time.Time `json:"dateRegistered"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Sale struct {
	ID         string
	CustomerID string
	SoldBy     string
	SaleDate   string
	TotalPrice float64
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
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
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
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
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// AccessoryMake represents the available accessory brands/makes
type AccessoryMake string

// AccessoryColor represents the available accessory colors
type AccessoryColor string

// AccessoryStatus represents the available inventory statuses
type AccessoryStatus string

// Constants for accessory makes
const (
	MakeGeneric     AccessoryMake = "Generic"
	MakeOEM         AccessoryMake = "OEM"
	MakeAftermarket AccessoryMake = "Aftermarket"
	MakeCustom      AccessoryMake = "Custom"
)

// Constants for accessory colors
const (
	ColorBlack  AccessoryColor = "Black"
	ColorWhite  AccessoryColor = "White"
	ColorSilver AccessoryColor = "Silver"
	ColorChrome AccessoryColor = "Chrome"
	ColorCustom AccessoryColor = "Custom"
)

// Constants for accessory statuses
const (
	StatusInStock    AccessoryStatus = "In Stock"
	StatusLowStock   AccessoryStatus = "Low Stock"
	StatusOutOfStock AccessoryStatus = "Out of Stock"
	StatusAvailable  AccessoryStatus = "Available"
)

// Accessory represents an accessory item in the inventory
type Accessory struct {
	ID        int             `json:"id"`         // Unique identifier
	Name      string          `json:"name"`       // Name of the accessory
	Make      AccessoryMake   `json:"make"`       // Manufacturer/brand of the accessory
	Quantity  int             `json:"quantity"`   // Number of units available
	Price     float64         `json:"price"`      // Price in PHP
	Status    AccessoryStatus `json:"status"`     // Inventory status
	UnitColor AccessoryColor  `json:"unit_color"` // Color of the accessory
	Image     string          `json:"image"`      // URL or base64 string of the image
	CreatedAt time.Time       `json:"createdAt"`  // Timestamp of creation
	UpdatedAt time.Time       `json:"updatedAt"`  // Timestamp of last update
}

// NewAccessoryInput represents data required to create a new accessory
type NewAccessoryInput struct {
	Name      string         `json:"name" validate:"required"`
	Make      AccessoryMake  `json:"make" validate:"required"`
	Quantity  int            `json:"quantity" validate:"required,min=0"`
	Price     float64        `json:"price" validate:"required,min=0"`
	UnitColor AccessoryColor `json:"unit_color" validate:"required"`
	Image     string         `json:"image"`
}

// UpdateAccessoryInput represents data to update an existing accessory
type UpdateAccessoryInput struct {
	Name      *string         `json:"name"`
	Make      *AccessoryMake  `json:"make"`
	Quantity  *int            `json:"quantity" validate:"omitempty,min=0"`
	Price     *float64        `json:"price" validate:"omitempty,min=0"`
	UnitColor *AccessoryColor `json:"unit_color"`
	Image     *string         `json:"image"`
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

// AccessoryForSale represents an accessory included in a cab sale
type AccessoryForSale struct {
	ID        int     `json:"id"`        // Accessory ID
	Name      string  `json:"name"`      // Name of the accessory
	Price     float64 `json:"price"`     // Price per unit
	Quantity  int     `json:"quantity"`  // Quantity being sold
	UnitPrice float64 `json:"unitPrice"` // Price per unit (same as Price)
}

// CabSalePayload represents the data sent from the frontend to record a cab sale
type CabSalePayload struct {
	CustomerID  string            `json:"customerId" validate:"required"` // ID of the customer making the purchase
	Quantity    int               `json:"quantity" validate:"required,min=1"` // Number of cabs being sold
	Accessories []AccessoryForSale `json:"accessories"` // Optional accessories included in the sale
}

// CabSale represents a completed cab sale transaction
type CabSale struct {
	CabID      int               `json:"cabId"`      // ID of the cab that was sold
	CustomerID string            `json:"customerId"` // ID of the customer who made the purchase
	Quantity   int               `json:"quantity"`   // Number of cabs sold
	Accessories []map[string]interface{} `json:"accessories"` // Accessories included in the sale
	TotalPrice float64           `json:"totalPrice"` // Total price of the sale
	SaleDate   string            `json:"saleDate"`   // Date of the sale
}
