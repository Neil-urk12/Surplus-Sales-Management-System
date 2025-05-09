package models

import (
	"time"
)

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
