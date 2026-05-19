package domain

import (
	"time"
)

type Category struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  *string   `json:"description,omitempty"`
	Tone         string    `json:"tone"`
	ImageURL     *string   `json:"image_url,omitempty"`
	SortOrder    int       `json:"sort_order"`
	IsActive     bool      `json:"is_active"`
	ProductCount int       `json:"product_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type VariantAttr struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type VariantOption struct {
	Key    string   `json:"key"`
	Label  string   `json:"label"`
	Values []string `json:"values"`
}

type Image struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	URL                string    `json:"url"`
	CloudinaryPublicID string    `json:"cloudinary_public_id"`
	CreatedAt          time.Time `json:"created_at"`
}

type Product struct {
	ID             string            `json:"id"`
	Title          string            `json:"title"`
	Subtitle       *string           `json:"subtitle,omitempty"`
	CategoryID     string            `json:"category_id"`
	Badge          *string           `json:"badge,omitempty"`
	BasePrice      int64             `json:"base_price"`
	Description    *string           `json:"description,omitempty"`
	Meaning        *string           `json:"meaning,omitempty"`
	VariantOptions []VariantOption   `json:"variant_options"`
	DefaultVariant map[string]string `json:"default_variant"`
	ZodiacIDs      []string          `json:"zodiac_ids"`
	PurposePlace   []string          `json:"purpose_place,omitempty"`
	PurposeUse     []string          `json:"purpose_use,omitempty"`
	PurposeAvoid   []string          `json:"purpose_avoid,omitempty"`
	Specs          map[string]string `json:"specs,omitempty"`
	RequiresSize   bool              `json:"requires_size"`
	IsActive       bool              `json:"is_active"`
	SortOrder      int               `json:"sort_order"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type ProductImage struct {
	ID        string        `json:"id"`
	ProductID string        `json:"product_id"`
	ImageID   string        `json:"image_id"`
	URL       string        `json:"url"`
	Name      string        `json:"name"`
	SortOrder int           `json:"sort_order"`
	Attrs     []VariantAttr `json:"attrs"`
	CreatedAt time.Time     `json:"created_at"`
}

type ProductSize struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	SizeLabel string `json:"size_label"`
	SizeCode  string `json:"size_code"`
	Price     int64  `json:"price"`
	SortOrder int    `json:"sort_order"`
}

type ProductSKU struct {
	ID        string            `json:"id"`
	ProductID string            `json:"product_id"`
	SizeCode  *string           `json:"size_code,omitempty"`
	Attrs     map[string]string `json:"attrs"`
	Price     int64             `json:"price"`
	SortOrder int               `json:"sort_order"`
	CreatedAt time.Time         `json:"created_at"`
}

type Campaign struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	DiscountType  string    `json:"discount_type"`
	DiscountValue int64     `json:"discount_value"`
	StartsAt      time.Time `json:"starts_at"`
	EndsAt        time.Time `json:"ends_at"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Review struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	ReviewerName string    `json:"reviewer_name"`
	Rating       int       `json:"rating"`
	Body         *string   `json:"body,omitempty"`
	IsApproved   bool      `json:"is_approved"`
	CreatedAt    time.Time `json:"created_at"`
}

type Banner struct {
	ID        string `json:"id"`
	Title     *string `json:"title,omitempty"`
	Subtitle  *string `json:"subtitle,omitempty"`
	ImageURL  *string `json:"image_url,omitempty"`
	LinkURL   *string `json:"link_url,omitempty"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ContactLink struct {
	ID        string    `json:"id"`
	Platform  string    `json:"platform"`
	Label     string    `json:"label"`
	URL       string    `json:"url"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID           string      `json:"id"`
	Phone        string      `json:"phone"`
	CustomerName *string     `json:"customer_name,omitempty"`
	Address      *string     `json:"address,omitempty"`
	Note         *string     `json:"note,omitempty"`
	Status       string      `json:"status"`
	AdminNote    *string     `json:"admin_note,omitempty"`
	TotalAmount  int64       `json:"total_amount"`
	Items        []OrderItem `json:"items"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID              string            `json:"id"`
	OrderID         string            `json:"order_id"`
	ProductID       string            `json:"product_id"`
	ProductTitle    string            `json:"product_title"`
	ProductSubtitle *string           `json:"product_subtitle,omitempty"`
	SizeCode        *string           `json:"size_code,omitempty"`
	SizeLabel       *string           `json:"size_label,omitempty"`
	SelectedAttrs   map[string]string `json:"selected_attrs,omitempty"`
	Quantity        int               `json:"quantity"`
	UnitPrice       int64             `json:"unit_price"`
	VariantImageURL *string           `json:"variant_image_url,omitempty"`
}

type WishlistItem struct {
	ID        string    `json:"id"`
	Phone     string    `json:"phone"`
	ProductID string    `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminUser struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	DisplayName  *string   `json:"display_name,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

type CustomerPhoto struct {
	ID        string    `json:"id"`
	ImageURL  string    `json:"image_url"`
	Caption   *string   `json:"caption,omitempty"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
