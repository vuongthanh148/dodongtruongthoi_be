package domain

import (
	"database/sql"
	"time"
)

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description *string `json:"description,omitempty"`
	Tone        string `json:"tone"`
	ImageURL    *string `json:"image_url,omitempty"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Product struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Subtitle      *string   `json:"subtitle,omitempty"`
	CategoryID    string    `json:"category_id"`
	Badge         *string   `json:"badge,omitempty"`
	BasePrice     int64     `json:"base_price"`
	Description   *string   `json:"description,omitempty"`
	Meaning       *string   `json:"meaning,omitempty"`
	DefaultBG     string            `json:"default_bg"`
	DefaultFrame  string            `json:"default_frame"`
	BGTones       []string          `json:"bg_tones"`
	Frames        []string          `json:"frames"`
	ZodiacIDs     []string          `json:"zodiac_ids"`
	PurposePlace  []string          `json:"purpose_place,omitempty"`
	PurposeUse    []string          `json:"purpose_use,omitempty"`
	PurposeAvoid  []string          `json:"purpose_avoid,omitempty"`
	Specs         map[string]string `json:"specs,omitempty"`
	RequiresBGTone bool             `json:"requires_bg_tone"`
	RequiresFrame bool              `json:"requires_frame"`
	RequiresSize  bool              `json:"requires_size"`
	IsActive      bool              `json:"is_active"`
	SortOrder     int               `json:"sort_order"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type ProductImage struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	BGTone    *string   `json:"bg_tone,omitempty"`
	Frame     *string   `json:"frame,omitempty"`
	URL       string    `json:"url"`
	AltText   *string   `json:"alt_text,omitempty"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

type ProductSize struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	SizeLabel string `json:"size_label"`
	SizeCode  string `json:"size_code"`
	Price     int64  `json:"price"`
	SortOrder int    `json:"sort_order"`
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
	Note         *string     `json:"note,omitempty"`
	Status       string      `json:"status"`
	AdminNote    *string     `json:"admin_note,omitempty"`
	TotalAmount  int64       `json:"total_amount"`
	Items        []OrderItem `json:"items"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID              string `json:"id"`
	OrderID         string `json:"order_id"`
	ProductID       string `json:"product_id"`
	ProductTitle    string `json:"product_title"`
	ProductSubtitle *string `json:"product_subtitle,omitempty"`
	SizeCode        *string `json:"size_code,omitempty"`
	SizeLabel       *string `json:"size_label,omitempty"`
	BGTone          *string `json:"bg_tone,omitempty"`
	BGToneLabel     *string `json:"bg_tone_label,omitempty"`
	Frame           *string `json:"frame,omitempty"`
	FrameLabel      *string `json:"frame_label,omitempty"`
	Quantity        int    `json:"quantity"`
	UnitPrice       int64  `json:"unit_price"`
	VariantImageURL *string `json:"variant_image_url,omitempty"`
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
	LastLoginAt  sql.NullTime `json:"last_login_at,omitempty"`
}
