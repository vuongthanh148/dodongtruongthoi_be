-- Seed Initial Data for Đồ Đồng Trường Thơi
-- Run with: psql -h localhost -U dodongtruongthoi -d dodongtruongthoi -f scripts/seed.sql

-- ===== Categories =====
INSERT INTO categories (id, name, slug, description, tone, image_url, sort_order, is_active)
VALUES
  ('tranh-phong-thuy', 'Tranh Phong Thủy', 'tranh-phong-thuy', 'Tranh đồng phong thủy truyền thống', 'gold', NULL, 1, true),
  ('dinh-dong-tho-cung', 'Đỉnh Đồng Thờ Cúng', 'dinh-dong-tho-cung', 'Đỉnh đồng cao cấp cho bàn thờ', 'bronze', NULL, 2, true),
  ('tuong-dong-trang-tri', 'Tượng Đồng Trang Trí', 'tuong-dong-trang-tri', 'Tượng đồng trang trí nội thất', 'bronze', NULL, 3, true),
  ('do-dung-nha-bep', 'Đồ Dùng Nhà Bếp', 'do-dung-nha-bep', 'Các đồ dùng nhà bếp bằng đồng', 'bronze', NULL, 4, true),
  ('phu-kien-trang-tri', 'Phụ Kiện Trang Trí', 'phu-kien-trang-tri', 'Phụ kiện trang trí nhà cửa', 'gold', NULL, 5, true)
ON CONFLICT (id) DO NOTHING;

-- ===== Products =====
INSERT INTO products (
  id, title, subtitle, category_id, badge, base_price, description, meaning,
  default_bg, default_frame, bg_tones, frames, zodiac_ids,
  purpose_place, purpose_use, purpose_avoid, specs,
  requires_bg_tone, requires_frame, requires_size, is_active, sort_order
)
VALUES
  (
    'tranh-nui-nuoc',
    'Tranh Núi Nước',
    'Tranh phong thủy truyền thống',
    'tranh-phong-thuy',
    'bestseller',
    2500000,
    'Tranh vẽ cảnh núi nước tuyệt đẹp, thể hiện sức mạnh và ổn định của thiên nhiên. Phù hợp với phong thủy của phòng khách và phòng làm việc.',
    'Biểu tượng của sự vĩnh cửu, ổn định và thịnh vượng. Năng lượng của núi nước giúp cân bằng không gian sống.',
    'gold', 'bronze',
    '{red,gold,bronze}', '{bronze,gold,dark}', '{tiger,dragon,horse}',
    '{living_room,office}', '{feng_shui_decor,housewarming_gift}', '{bedroom}', '{"material":"Đồng","technique":"Chạm nổi"}'::jsonb,
    true, true, true, true, 1
  ),
  (
    'tranh-hoa-tim',
    'Tranh Hoa Tím',
    'Tranh phong thủy hoa tím cao quý',
    'tranh-phong-thuy',
    'new',
    1800000,
    'Tranh đồng chạm nổi hoa tím biểu tượng của sự sang trọng và quý phái. Hoa tím mang lại bình yên và tình yêu.',
    'Hoa tím đại diện cho sự tinh tế, bình yên và tình yêu vợ chồng. Phù hợp với phòng ngủ hay phòng khách.',
    'gold', 'gold',
    '{gold,bronze}', '{gold,dark}', '{rabbit,goat}',
    '{living_room,bedroom}', '{housewarming_gift,anniversary_gift}', '{}', '{"material":"Đồng","size":"1.2m × 0.8m"}'::jsonb,
    true, true, true, true, 2
  ),
  (
    'dinh-dong-3chan',
    'Đỉnh Đồng 3 Chân',
    'Đỉnh đồng thờ cúng truyền thống',
    'dinh-dong-tho-cung',
    NULL,
    4500000,
    'Đỉnh đồng 3 chân kiểu cổ, dùng để thờ cúng tổ tiên và các vị thần linh. Chế tác tinh xảo, đỏm dáng thanh lịch.',
    'Đỉnh đồng 3 chân là vật phẩm thiêng liêng trong văn hóa tâm linh, giúp kết nối với tổ tiên và vũ trụ.',
    'bronze', 'bronze',
    '{bronze,gold}', '{bronze}', '{}',
    '{shrine,altar}', '{}', '{}', '{"height":"0.5m","weight":"15kg","material":"Đồng nguyên chất"}'::jsonb,
    true, false, true, true, 3
  ),
  (
    'tuong-phat-a-di-da',
    'Tượng Phật A Di Đà',
    'Tượng Phật mang lại bình yên',
    'tuong-dong-trang-tri',
    'sale',
    3200000,
    'Tượng Phật A Di Đà được chạm khắc tỉ mỉ, biểu tượng của sự giáng phúc và bình yên. Có thể để ở bàn thờ hay trang trí nhà cửa.',
    'Phật A Di Đà mang lại bình yên, may mắn và bảo vệ toàn gia đình. Là vật phẩm tâm linh cao quý.',
    'bronze', 'bronze',
    '{bronze}', '{bronze,carved}', '{}',
    '{shrine,living_room}', '{spiritual_protection}', '{}', '{"height":"0.6m","weight":"8kg"}'::jsonb,
    false, true, false, true, 4
  ),
  (
    'nhu-y-dong',
    'Như Ý Đồng',
    'Như ý đồng phong thủy cao cấp',
    'phu-kien-trang-tri',
    NULL,
    1200000,
    'Như ý đồng là biểu tượng của sự may mắn, thành công và ổn định. Thích hợp để làm quà tặng hoặc trang trí bàn làm việc.',
    'Như ý đưa đến sự bình yên, thành công trong công việc và những điều tốt đẹp trong cuộc sống.',
    'gold', 'gold',
    '{gold,red}', '{gold}', '{}',
    '{office,living_room}', '{good_luck,career_success}', '{}', '{"material":"Đồng","height":"0.3m"}'::jsonb,
    true, true, false, true, 5
  )
ON CONFLICT (id) DO NOTHING;

-- ===== Product Sizes =====
INSERT INTO product_sizes (product_id, size_label, size_code, price, sort_order)
VALUES
  ('tranh-nui-nuoc', '0.8m × 0.6m', 's', 2500000, 1),
  ('tranh-nui-nuoc', '1.2m × 0.8m', 'm', 3500000, 2),
  ('tranh-nui-nuoc', '1.5m × 1.0m', 'l', 4800000, 3),
  ('tranh-nui-nuoc', '2.0m × 1.3m', 'xl', 6500000, 4),
  ('tranh-hoa-tim', '0.8m × 0.6m', 's', 1800000, 1),
  ('tranh-hoa-tim', '1.2m × 0.8m', 'm', 2400000, 2),
  ('tranh-hoa-tim', '1.5m × 1.0m', 'l', 3200000, 3),
  ('dinh-dong-3chan', 'Tiêu chuẩn', 's', 4500000, 1),
  ('dinh-dong-3chan', 'Lớn', 'l', 6500000, 2),
  ('tuong-phat-a-di-da', 'Tiêu chuẩn', 's', 3200000, 1),
  ('tuong-phat-a-di-da', 'Lớn', 'l', 4200000, 2),
  ('nhu-y-dong', 'Tiêu chuẩn', 's', 1200000, 1)
ON CONFLICT (product_id, size_code) DO NOTHING;

-- ===== Product Images =====
INSERT INTO product_images (product_id, bg_tone, frame, url, alt_text, sort_order)
VALUES
  ('tranh-nui-nuoc', 'gold', 'bronze', 'https://res.cloudinary.com/dodongtruongthoi/image/upload/v1777016322/products/skqilesdnltbpg99btho.jpg', 'Tranh Núi Nước 1', 1),
  ('tranh-nui-nuoc', 'gold', 'bronze', 'https://res.cloudinary.com/dodongtruongthoi/image/upload/v1777016322/products/skqilesdnltbpg99btho.jpg', 'Tranh Núi Nước 2', 2),
  ('tranh-nui-nuoc', 'gold', 'bronze', 'https://res.cloudinary.com/dodongtruongthoi/image/upload/v1777016322/products/skqilesdnltbpg99btho.jpg', 'Tranh Núi Nước 3', 3);

-- ===== Banners =====
INSERT INTO banners (title, subtitle, image_url, link_url, sort_order, is_active)
VALUES
  ('Tết 2026 - Ưu Đãi Đặc Biệt', 'Giảm giá 30% cho tất cả sản phẩm', NULL, '/products?campaign=tet-2026', 1, true),
  ('Về Chúng Tôi', 'Thủ công mỹ nghệ truyền thống 50 năm', NULL, '/about', 2, true),
  ('Chất Lượng Tuyệt Vời', 'Mỗi sản phẩm đều được chạm khắc tỉ mỉ', NULL, '/collections/bestsellers', 3, true)
ON CONFLICT DO NOTHING;

-- ===== Contact Links =====
INSERT INTO contact_links (platform, label, url, sort_order, is_active)
VALUES
  ('zalo', 'Zalo', 'https://zalo.me/0899012288', 1, true),
  ('messenger', 'Messenger', 'https://m.me/dodongtruongthoi', 2, true),
  ('facebook', 'Facebook', 'https://facebook.com/dodongtruongthoi', 3, true),
  ('tiktok', 'TikTok', 'https://tiktok.com/@dodongtruongthoi', 4, true),
  ('phone', 'Hotline', 'tel:0899012288', 5, true),
  ('email', 'Email', 'mailto:info@dodongtruongthoi.com', 6, true)
ON CONFLICT DO NOTHING;

-- ===== Campaigns =====
INSERT INTO campaigns (name, description, discount_type, discount_value, starts_at, ends_at, is_active)
VALUES
  ('Tết 2026 - Khuyến Mãi', 'Giảm 30% cho tất cả tranh phong thủy trong dịp Tết', 'percentage', 30, NOW(), NOW() + INTERVAL '30 days', true),
  ('Khuyến Mãi Mùa Hè', 'Giảm 20% cho các tượng và đỉnh đồng', 'percentage', 20, NOW() + INTERVAL '60 days', NOW() + INTERVAL '90 days', false)
ON CONFLICT DO NOTHING;

-- ===== Campaign Products =====
INSERT INTO campaign_products (campaign_id, product_id)
SELECT c.id, 'tranh-nui-nuoc'::text
FROM campaigns c WHERE c.name = 'Tết 2026 - Khuyến Mãi'
UNION ALL
SELECT c.id, 'tranh-hoa-tim'::text
FROM campaigns c WHERE c.name = 'Tết 2026 - Khuyến Mãi'
UNION ALL
SELECT c.id, 'tuong-phat-a-di-da'::text
FROM campaigns c WHERE c.name = 'Khuyến Mãi Mùa Hè'
UNION ALL
SELECT c.id, 'dinh-dong-3chan'::text
FROM campaigns c WHERE c.name = 'Khuyến Mãi Mùa Hè'
ON CONFLICT (campaign_id, product_id) DO NOTHING;

-- ===== Reviews =====
INSERT INTO reviews (product_id, reviewer_name, rating, body, is_approved)
VALUES
  ('tranh-nui-nuoc', 'Anh Hùng', 5, 'Tranh đẹp lắm, chạm khắc tinh xảo. Giao hàng nhanh và an toàn.', true),
  ('tranh-nui-nuoc', 'Chị Linh', 5, 'Mình mua để treo phòng khách, rất hài lòng. Sẽ mua lại.', true),
  ('tranh-nui-nuoc', 'Anh Minh', 4, 'Sản phẩm tốt, nhưng giá hơi cao.', true),
  ('tranh-hoa-tim', 'Chị Thu', 5, 'Đẹp quá! Phòng ngủ của mình trở nên sang trọng hơn.', true),
  ('dinh-dong-3chan', 'Anh Tâm', 5, 'Đỉnh đồng chất lượng cao, rất đẹp.', true),
  ('nhu-y-dong', 'Chị Hoa', 4, 'Tặng chồng làm quà, anh ấy rất thích.', true)
ON CONFLICT DO NOTHING;

-- ===== Summary =====
\echo 'Seed data loaded successfully!'
SELECT 'Categories' as entity, COUNT(*) as count FROM categories
UNION ALL
SELECT 'Products', COUNT(*) FROM products
UNION ALL
SELECT 'Product Sizes', COUNT(*) FROM product_sizes
UNION ALL
SELECT 'Banners', COUNT(*) FROM banners
UNION ALL
SELECT 'Contacts', COUNT(*) FROM contact_links
UNION ALL
SELECT 'Campaigns', COUNT(*) FROM campaigns
UNION ALL
SELECT 'Reviews', COUNT(*) FROM reviews
ORDER BY entity;
