# Test Data Dictionary

## Overview

This document describes all test data created for comprehensive admin CMS testing. Data includes realistic Vietnamese names, phone numbers, and product information from the Đồ Đồng Trường Thơi catalog.

---

## Test Orders (14 total)

All orders contain real product combinations and realistic customer information.

### By Status

#### Pending Confirmation (3 orders)
| ID | Phone | Customer | Amount | Items | Note |
|----|-------|----------|--------|-------|------|
| ... | 0901234567 | Nguyễn Văn A | 7,000,000 VND | 2 | "Giao hàng tại nhà số 123 Nguyễn Huệ" |
| ... | 0912345678 | Trần Thị B | 4,500,000 VND | 1 | "Giao hàng tối" |
| ... | 0923456789 | Lê Văn C | 9,200,000 VND | 2 | "Vui lòng gọi trước khi giao" |

**Test Case:** Filter by status, display pending orders for confirmation

#### Confirmed (2 orders)
| ID | Phone | Customer | Amount | Items | Admin Note |
|----|-------|----------|--------|-------|------------|
| ... | 0934567890 | Phạm Thị D | 5,500,000 VND | 2 | "Đã liên hệ khách hàng, xác nhận địa chỉ" |
| ... | 0945678901 | Hoàng Văn E | 3,200,000 VND | 1 | "Khách hàng đã thanh toán cọc" |

**Test Case:** Orders ready for processing

#### Processing (3 orders)
| ID | Phone | Customer | Amount | Items | Admin Note |
|----|-------|----------|--------|-------|------------|
| ... | 0956789012 | Vũ Thị F | 6,800,000 VND | 1 | "Đang chuẩn bị sản phẩm" |
| ... | 0967890123 | Bùi Văn G | 4,200,000 VND | 2 | "Đang đóng gói, sẵn sàng giao" |

**Test Case:** Orders being prepared for shipment

#### Shipped (2 orders)
| ID | Phone | Customer | Amount | Items | Admin Note |
|----|-------|----------|--------|-------|------------|
| ... | 0978901234 | Cao Thị H | 7,500,000 VND | 2 | "Giao hàng thứ hai" |
| ... | 0989012345 | Đặng Văn I | 5,200,000 VND | 1 | "Đơn hàng đã gửi đi" |

**Test Case:** Orders in transit, visible in customer tracking

#### Completed (2 orders)
| ID | Phone | Customer | Amount | Items | Admin Note |
|----|-------|----------|--------|-------|------------|
| ... | 0990123456 | Tôn Thị J | 8,900,000 VND | 2 | "Giao hàng thành công, khách hàng xác nhận" |
| ... | 0901234500 | Võ Văn K | 3,800,000 VND | 1 | "Hoàn tất đơn hàng" |

**Test Case:** Orders successfully delivered, appear in order history

#### Cancelled (2 orders)
| ID | Phone | Customer | Amount | Items | Reason |
|----|-------|----------|--------|-------|--------|
| ... | 0912345500 | Nông Thị L | 6,200,000 VND | 1 | "Khách hàng hủy vì lý do cá nhân" |
| ... | 0923456500 | Mạc Văn M | 4,900,000 VND | 2 | "Không liên lạc được với khách hàng" |

**Test Case:** Handling cancellations and order refunds

---

## Test Order Items (19 total)

Orders contain realistic product combinations with full variant details.

### Product Mix

| Product | Orders Used | Frequency |
|---------|-------------|-----------|
| Tranh Núi Nước | 6 | Most popular |
| Tranh Hoa Tím | 4 | Popular |
| Nhu Ý Đồng | 3 | Regular |
| Đỉnh Đồng 3 Chân | 3 | Popular |
| Tượng Phật A Di Đà | 3 | Popular |

### Variant Coverage

**Sizes Tested:**
- s: "0.8m × 0.6m"
- m: "1.2m × 0.8m"
- l: "1.5m × 1.0m"
- xl: "2.0m × 1.3m"
- Standard (for non-sized items)

**Colors Tested:**
- gold (Vàng)
- bronze (Đồng)
- red (Đỏ)

**Frames Tested:**
- bronze
- gold
- dark
- carved

**Sample Item:**
```
Order ID: ...
Product: Tranh Núi Nước
Size: m (1.2m × 0.8m)
Color: gold (Vàng)
Frame: bronze (Đồng)
Quantity: 1
Unit Price: 3,500,000 VND
Total: 3,500,000 VND
```

---

## Test Reviews (19 total)

### Distribution

| Status | Count | % | Average Rating |
|--------|-------|---|-----------------|
| Approved | 13 | 68% | 4.6/5 |
| Pending | 6 | 32% | 4.3/5 |
| **Total** | **19** | **100%** | **4.5/5** |

### By Product

| Product | Approved | Pending | Total |
|---------|----------|---------|-------|
| Tranh Núi Nước | 4 | 1 | 5 |
| Tranh Hoa Tím | 3 | 1 | 4 |
| Tượng Phật A Di Đà | 2 | 1 | 3 |
| Nhu Ý Đồng | 2 | 1 | 3 |
| Dinh Đồng 3 Chân | 1 | 1 | 2 |
| Other Products | 1 | 1 | 2 |

### Rating Distribution

| Rating | Count | Examples |
|--------|-------|----------|
| 5★ | 11 | "Tuyệt vời!", "Đẹp lắm!" |
| 4★ | 5 | "Tốt, nhưng...", "Hài lòng" |
| 3★ | 1 | "Bình thường" |

### Sample Approved Review
```
Reviewer: Anh Hùng
Product: Tranh Núi Nước
Rating: 5★
Comment: "Tranh đẹp lắm, chạm khắc tinh xảo. Giao hàng nhanh và an toàn."
Status: Approved
Created: [Past date]
```

### Sample Pending Review
```
Reviewer: Anh Mạnh
Product: Tranh Núi Nước
Rating: 5★
Comment: "Tuyệt vời! Chạm khắc rất tinh xảo, giao hàng nhanh."
Status: Pending
Created: [Recent date]
```

---

## Test Campaigns (2 total)

### Campaign 1: Tết 2026 - Khuyến Mãi (Active)
```
ID: [UUID]
Name: Tết 2026 - Khuyến Mãi
Description: Giảm 30% cho tất cả tranh phong thủy trong dịp Tết
Discount Type: percentage
Discount Value: 30%
Starts At: NOW()
Ends At: NOW() + 30 days
Active: Yes
Products Assigned: tranh-nui-nuoc, tranh-hoa-tim, tuong-phat-a-di-da, dinh-dong-3chan
```

**Test Case:** Active campaign with products, applicable to checkout

### Campaign 2: Khuyến Mãi Mùa Hè (Inactive)
```
ID: [UUID]
Name: Khuyến Mãi Mùa Hè
Description: Giảm 20% cho các tượng và đỉnh đồng
Discount Type: percentage
Discount Value: 20%
Starts At: NOW() + 60 days
Ends At: NOW() + 90 days
Active: No
Products Assigned: tuong-phat-a-di-da, dinh-dong-3chan
```

**Test Case:** Future campaign not yet active

---

## Existing Products (5 total)

All products exist and are used in orders and campaigns.

| Product ID | Title | Category | Price | Active |
|------------|-------|----------|-------|--------|
| tranh-nui-nuoc | Tranh Núi Nước | Tranh Phong Thủy | 2,500,000 | Yes |
| tranh-hoa-tim | Tranh Hoa Tím | Tranh Phong Thủy | 1,800,000 | Yes |
| dinh-dong-3chan | Đỉnh Đồng 3 Chân | Đỉnh Đồng Thờ Cúng | 4,500,000 | Yes |
| tuong-phat-a-di-da | Tượng Phật A Di Đà | Tượng Đồng Trang Trí | 3,200,000 | Yes |
| nhu-y-dong | Như Ý Đồng | Phụ Kiện Trang Trí | 1,200,000 | Yes |

---

## Existing Categories (5 total)

| Category ID | Name | Description | Active |
|------------|------|-------------|--------|
| tranh-phong-thuy | Tranh Phong Thủy | Tranh đồng phong thủy truyền thống | Yes |
| dinh-dong-tho-cung | Đỉnh Đồng Thờ Cúng | Đỉnh đồng cao cấp cho bàn thờ | Yes |
| tuong-dong-trang-tri | Tượng Đồng Trang Trí | Tượng đồng trang trí nội thất | Yes |
| do-dung-nha-bep | Đồ Dùng Nhà Bếp | Các đồ dùng nhà bếp bằng đồng | Yes |
| phu-kien-trang-tri | Phụ Kiện Trang Trí | Phụ kiện trang trí nhà cửa | Yes |

---

## Test Banners (6 total)

### Original (3)
```
1. Tết 2026 - Ưu Đãi Đặc Biệt
   Subtitle: Giảm giá 30% cho tất cả sản phẩm
   Link: /products?campaign=tet-2026
   Active: Yes

2. Về Chúng Tôi
   Subtitle: Thủ công mỹ nghệ truyền thống 50 năm
   Link: /about
   Active: Yes

3. Chất Lượng Tuyệt Vời
   Subtitle: Mỗi sản phẩm đều được chạm khắc tỉ mỉ
   Link: /collections/bestsellers
   Active: Yes
```

### Test Data (3)
```
4. Khuyến Mãi Black Friday
   Subtitle: Giảm giá lên đến 50% cho các sản phẩm chọn lọc
   Link: /products?campaign=black-friday
   Active: No (scheduled for future)

5. Nơi Tạo Ra Chất Lượng
   Subtitle: Mỗi sản phẩm là kết quả của tâm huyết và tư tưởng
   Link: /about
   Active: Yes

6. Giao Hàng Miễn Phí
   Subtitle: Với mọi đơn hàng trên 5 triệu đồng
   Link: /shop
   Active: Yes
```

**Test Cases:**
- Active banners display in carousel
- Inactive banners hidden from storefront
- Banners are sorted by sort_order
- Banner clicks navigate correctly

---

## Contact Links (6 total)

Represents all communication channels.

| Platform | Label | URL | Active |
|----------|-------|-----|--------|
| zalo | Zalo | https://zalo.me/0899012288 | Yes |
| messenger | Messenger | https://m.me/dodongtruongthoi | Yes |
| facebook | Facebook | https://facebook.com/dodongtruongthoi | Yes |
| tiktok | TikTok | https://tiktok.com/@dodongtruongthoi | Yes |
| phone | Hotline | tel:0899012288 | Yes |
| email | Email | mailto:info@dodongtruongthoi.com | Yes |

---

## Data Characteristics

### Realistic Details
- ✅ Vietnamese names with proper formatting
- ✅ Valid Vietnamese phone numbers (starting with 09)
- ✅ Realistic customer notes and shipping instructions
- ✅ Actual products from catalog
- ✅ Believable price points and order values

### Edge Cases Covered
- ✅ Orders with 1 item
- ✅ Orders with multiple items (up to 3)
- ✅ Multiple variants of same product
- ✅ Orders with and without customer names
- ✅ Orders with and without admin notes
- ✅ All order statuses represented

### Quantity Coverage
- ✅ Single-item orders
- ✅ Multi-item orders
- ✅ Orders with quantity > 1 per item
- ✅ Various product combinations

---

## How to Use This Data

### For Manual Testing
1. Login to admin at http://localhost:3000/admin/login
2. Use credentials: admin / admin123
3. Navigate to each module and verify data displays correctly

### For API Testing
1. Get auth token: POST /admin/login
2. Use token in Authorization header
3. Call endpoints like:
   - GET /admin/orders - returns 14 orders
   - GET /admin/orders?status=pending_confirm - returns 3 orders
   - GET /admin/reviews?approved=false - returns 6 pending reviews

### For Performance Testing
- Current dataset is suitable for basic load testing
- Scale by running seed script multiple times
- For 100+ orders, recommend pagination implementation

### For Feature Development
- Use existing data patterns as templates
- All relationships are properly set up
- Safe to create/modify/delete test data freely
- Re-run seed script to reset to baseline

---

## Data Reset

To reset to original seed data only (no test data):
```bash
# In backend directory
psql -h localhost -U dodongtruongthoi -d dodongtruongthoi -c "
TRUNCATE orders CASCADE;
TRUNCATE reviews CASCADE;
TRUNCATE campaigns CASCADE;
TRUNCATE campaign_products CASCADE;
TRUNCATE banners CASCADE;
"
```

To reload all test data:
```bash
psql -h localhost -U dodongtruongthoi -d dodongtruongthoi -f scripts/seed-test-data.sql
```

---

## Statistics

### Total Test Records
- Orders: 14
- Order Items: 19
- Reviews: 19 (13 approved, 6 pending)
- Campaigns: 2 (1 active, 1 future)
- Banners: 6 (5 active, 1 inactive)
- Products: 5
- Categories: 5
- Contacts: 6

### Financial Summary
- Total Revenue (Orders): 79,400,000 VND
- Average Order Value: 5,671,429 VND
- Price Range: 2,500,000 - 9,200,000 VND
- Total Items Ordered: 19

### Content Metrics
- Product Variants: 20+ combinations
- Review Coverage: All products reviewed (3-5 reviews each)
- Campaign Product Assignment: 4 products per campaign
- Contact Channels: All 6 platforms represented

---

**Last Updated:** April 25, 2026
**Seed Files:**
- `/scripts/seed.sql` - Original data
- `/scripts/seed-test-data.sql` - Test data additions
