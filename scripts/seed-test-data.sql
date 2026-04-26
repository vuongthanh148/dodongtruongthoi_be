-- Comprehensive Test Data for Admin CMS
-- Run with: psql -h localhost -U dodongtruongthoi -d dodongtruongthoi -f scripts/seed-test-data.sql

-- ===== Test Orders =====
-- Create diverse orders with different statuses and amounts

INSERT INTO orders (id, phone, customer_name, note, status, admin_note, total_amount, created_at)
VALUES
  -- Pending confirmation orders
  (gen_random_uuid(), '0901234567', 'Nguyễn Văn A', 'Giao hàng tại nhà số 123 Nguyễn Huệ', 'pending_confirm', NULL, 7000000, NOW() - INTERVAL '2 days'),
  (gen_random_uuid(), '0912345678', 'Trần Thị B', 'Giao hàng tối', 'pending_confirm', NULL, 4500000, NOW() - INTERVAL '1 day'),
  (gen_random_uuid(), '0923456789', 'Lê Văn C', 'Vui lòng gọi trước khi giao', 'pending_confirm', NULL, 9200000, NOW() - INTERVAL '3 hours'),

  -- Confirmed orders
  (gen_random_uuid(), '0934567890', 'Phạm Thị D', 'Giao hàng cuối tuần', 'confirmed', 'Đã liên hệ khách hàng, xác nhận địa chỉ', 5500000, NOW() - INTERVAL '1 day'),
  (gen_random_uuid(), '0945678901', 'Hoàng Văn E', NULL, 'confirmed', 'Khách hàng đã thanh toán cọc', 3200000, NOW() - INTERVAL '2 days'),

  -- Processing orders
  (gen_random_uuid(), '0956789012', 'Vũ Thị F', 'Giao hàng sáng sớm', 'processing', 'Đang chuẩn bị sản phẩm', 6800000, NOW() - INTERVAL '2 days'),
  (gen_random_uuid(), '0967890123', 'Bùi Văn G', 'Chuyển khách hàng khác', 'processing', 'Đang đóng gói, sẵn sàng giao', 4200000, NOW() - INTERVAL '1 day'),

  -- Shipped orders
  (gen_random_uuid(), '0978901234', 'Cao Thị H', 'Giao hàng qua nhân viên giao hàng', 'shipped', 'Giao hàng thứ hai', 7500000, NOW() - INTERVAL '3 days'),
  (gen_random_uuid(), '0989012345', 'Đặng Văn I', NULL, 'shipped', 'Đơn hàng đã gửi đi', 5200000, NOW() - INTERVAL '2 days'),

  -- Completed orders
  (gen_random_uuid(), '0990123456', 'Tôn Thị J', 'Khách hàng rất hài lòng', 'completed', 'Giao hàng thành công, khách hàng xác nhận', 8900000, NOW() - INTERVAL '5 days'),
  (gen_random_uuid(), '0901234500', 'Võ Văn K', NULL, 'completed', 'Hoàn tất đơn hàng', 3800000, NOW() - INTERVAL '4 days'),

  -- Cancelled orders
  (gen_random_uuid(), '0912345500', 'Nông Thị L', 'Khách hàng hủy vì lý do cá nhân', 'cancelled', 'Đã hoàn tiền 100%', 6200000, NOW() - INTERVAL '3 days'),
  (gen_random_uuid(), '0923456500', 'Mạc Văn M', 'Không liên lạc được với khách hàng', 'cancelled', 'Khách hàng không nhận điện thoại', 4900000, NOW() - INTERVAL '2 days');

-- ===== Insert Order Items for Each Order =====
-- Get order IDs and insert items
INSERT INTO order_items (order_id, product_id, product_title, product_subtitle, size_code, size_label, bg_tone, bg_tone_label, frame, frame_label, quantity, unit_price, variant_image_url)
SELECT
  (SELECT id FROM orders WHERE phone = '0901234567' LIMIT 1),
  'tranh-nui-nuoc', 'Tranh Núi Nước', 'Tranh phong thủy truyền thống',
  'm', '1.2m × 0.8m', 'gold', 'Vàng', 'bronze', 'Đồng', 1, 3500000,
  'https://res.cloudinary.com/dodongtruongthoi/image/upload/v1777016322/products/skqilesdnltbpg99btho.jpg'
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0901234567' LIMIT 1),
  'nhu-y-dong', 'Như Ý Đồng', 'Như ý đồng phong thủy cao cấp',
  's', 'Tiêu chuẩn', 'gold', 'Vàng', 'gold', 'Vàng', 1, 1200000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0912345678' LIMIT 1),
  'dinh-dong-3chan', 'Đỉnh Đồng 3 Chân', 'Đỉnh đồng thờ cúng truyền thống',
  'l', 'Lớn', 'bronze', 'Đồng', 'bronze', 'Đồng', 1, 6500000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0923456789' LIMIT 1),
  'tranh-hoa-tim', 'Tranh Hoa Tím', 'Tranh phong thủy hoa tím cao quý',
  'l', '1.5m × 1.0m', 'gold', 'Vàng', 'gold', 'Vàng', 1, 3200000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0923456789' LIMIT 1),
  'tuong-phat-a-di-da', 'Tượng Phật A Di Đà', 'Tượng Phật mang lại bình yên',
  's', 'Tiêu chuẩn', 'bronze', 'Đồng', 'bronze', 'Đồng', 2, 3200000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0934567890' LIMIT 1),
  'tranh-nui-nuoc', 'Tranh Núi Nước', 'Tranh phong thủy truyền thống',
  's', '0.8m × 0.6m', 'red', 'Đỏ', 'bronze', 'Đồng', 1, 2500000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0934567890' LIMIT 1),
  'tranh-hoa-tim', 'Tranh Hoa Tím', 'Tranh phong thủy hoa tím cao quý',
  's', '0.8m × 0.6m', 'bronze', 'Đồng', 'dark', 'Tối', 1, 1800000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0945678901' LIMIT 1),
  'tuong-phat-a-di-da', 'Tượng Phật A Di Đà', 'Tượng Phật mang lại bình yên',
  'l', 'Lớn', 'bronze', 'Đồng', 'carved', 'Chạm Khắc', 1, 4200000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0956789012' LIMIT 1),
  'tranh-nui-nuoc', 'Tranh Núi Nước', 'Tranh phong thủy truyền thống',
  'xl', '2.0m × 1.3m', 'gold', 'Vàng', 'gold', 'Vàng', 1, 6500000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0967890123' LIMIT 1),
  'nhu-y-dong', 'Như Ý Đồng', 'Như ý đồng phong thủy cao cấp',
  's', 'Tiêu chuẩn', 'gold', 'Vàng', 'gold', 'Vàng', 2, 1200000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0978901234' LIMIT 1),
  'dinh-dong-3chan', 'Đỉnh Đồng 3 Chân', 'Đỉnh đồng thờ cúng truyền thống',
  's', 'Tiêu chuẩn', 'bronze', 'Đồng', 'bronze', 'Đồng', 1, 4500000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0978901234' LIMIT 1),
  'tranh-hoa-tim', 'Tranh Hoa Tím', 'Tranh phong thủy hoa tím cao quý',
  'm', '1.2m × 0.8m', 'gold', 'Vàng', 'gold', 'Vàng', 1, 2400000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0989012345' LIMIT 1),
  'tranh-nui-nuoc', 'Tranh Núi Nước', 'Tranh phong thủy truyền thống',
  'm', '1.2m × 0.8m', 'bronze', 'Đồng', 'bronze', 'Đồng', 1, 3500000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0990123456' LIMIT 1),
  'tranh-nui-nuoc', 'Tranh Núi Nước', 'Tranh phong thủy truyền thống',
  'l', '1.5m × 1.0m', 'gold', 'Vàng', 'bronze', 'Đồng', 1, 4800000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0990123456' LIMIT 1),
  'tuong-phat-a-di-da', 'Tượng Phật A Di Đà', 'Tượng Phật mang lại bình yên',
  'l', 'Lớn', 'bronze', 'Đồng', 'bronze', 'Đồng', 1, 4200000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0901234500' LIMIT 1),
  'tranh-hoa-tim', 'Tranh Hoa Tím', 'Tranh phong thủy hoa tím cao quý',
  's', '0.8m × 0.6m', 'gold', 'Vàng', 'gold', 'Vàng', 1, 1800000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0912345500' LIMIT 1),
  'nhu-y-dong', 'Như Ý Đồng', 'Như ý đồng phong thủy cao cấp',
  's', 'Tiêu chuẩn', 'red', 'Đỏ', 'gold', 'Vàng', 3, 1200000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0923456500' LIMIT 1),
  'dinh-dong-3chan', 'Đỉnh Đồng 3 Chân', 'Đỉnh đồng thờ cúng truyền thống',
  'l', 'Lớn', 'bronze', 'Đồng', 'bronze', 'Đồng', 1, 6500000,
  NULL
UNION ALL
SELECT
  (SELECT id FROM orders WHERE phone = '0923456500' LIMIT 1),
  'tranh-nui-nuoc', 'Tranh Núi Nước', 'Tranh phong thủy truyền thống',
  's', '0.8m × 0.6m', 'gold', 'Vàng', 'bronze', 'Đồng', 1, 2500000,
  NULL;

-- ===== Additional Reviews (Pending & Approved) =====
INSERT INTO reviews (product_id, reviewer_name, rating, body, is_approved)
VALUES
  -- Pending reviews
  ('tranh-nui-nuoc', 'Anh Mạnh', 5, 'Tuyệt vời! Chạm khắc rất tinh xảo, giao hàng nhanh.', false),
  ('tranh-hoa-tim', 'Chị Mỹ', 4, 'Hình ảnh sản phẩm rất đẹp, nhưng màu sắc có hơi khác so với ảnh.', false),
  ('dinh-dong-3chan', 'Anh Tuyên', 3, 'Sản phẩm bình thường, không hơn gì lắm.', false),
  ('tuong-phat-a-di-da', 'Chị Hương', 5, 'Dùng để thờ cúng, rất linh thiêng. Tặng cho mẹ, mẹ rất vui.', false),
  ('nhu-y-dong', 'Anh Sáng', 5, 'Để bàn làm việc, mang lại may mắn cho công việc.', false),

  -- More approved reviews for variety
  ('tranh-nui-nuoc', 'Chị Trang', 5, 'Lần thứ 3 mua hàng của cửa hàng, lần nào cũng hài lòng.', true),
  ('tranh-hoa-tim', 'Anh Phong', 4, 'Chất lượng tốt, giá cả hợp lý. Sẽ giới thiệu bạn bè.', true),
  ('tuong-phat-a-di-da', 'Chị Hà', 5, 'Giao hàng nhanh, sản phẩm đúng như mô tả. Très magnifique!', true),
  ('dinh-dong-3chan', 'Anh Dũng', 5, 'Đỉnh đồng chất lượng cao, giá hơi cao nhưng đáng tiền.', true),
  ('nhu-y-dong', 'Chị Loan', 4, 'Sản phẩm đẹp, mang lại cảm giác bình yên.', true),
  ('tranh-nui-nuoc', 'Anh Tuấn', 5, 'Trang trí phòng khách, bạn bè ai vào đều khen.', true),
  ('tranh-hoa-tim', 'Chị Bích', 5, 'Phòng ngủ trở nên sang trọng gấp bội. Cảm ơn cửa hàng.', true);

-- ===== Additional Banners =====
INSERT INTO banners (title, subtitle, image_url, link_url, sort_order, is_active)
VALUES
  ('Khuyến Mãi Black Friday', 'Giảm giá lên đến 50% cho các sản phẩm chọn lọc', NULL, '/products?campaign=black-friday', 4, false),
  ('Nơi Tạo Ra Chất Lượng', 'Mỗi sản phẩm là kết quả của tâm huyết và tư tưởng', NULL, '/about', 5, true),
  ('Giao Hàng Miễn Phí', 'Với mọi đơn hàng trên 5 triệu đồng', NULL, '/shop', 6, true)
ON CONFLICT DO NOTHING;

-- ===== Summary Statistics =====
\echo '=========================================='
\echo 'Test Data Loaded Successfully!'
\echo '=========================================='
SELECT 'Orders' as entity, COUNT(*) as count FROM orders
UNION ALL
SELECT 'Order Items', COUNT(*) FROM order_items
UNION ALL
SELECT 'Reviews (Total)', COUNT(*) FROM reviews
UNION ALL
SELECT 'Reviews (Approved)', COUNT(*) FROM reviews WHERE is_approved = true
UNION ALL
SELECT 'Reviews (Pending)', COUNT(*) FROM reviews WHERE is_approved = false
UNION ALL
SELECT 'Banners', COUNT(*) FROM banners
ORDER BY entity;

-- ===== Order Status Distribution =====
\echo ''
\echo 'Order Status Distribution:'
\echo '------------------------------------------'
SELECT status, COUNT(*) as count FROM orders GROUP BY status ORDER BY status;

-- ===== Average Order Value =====
\echo ''
\echo 'Order Statistics:'
\echo '------------------------------------------'
SELECT
  'Average Order Value' as metric,
  ROUND(AVG(total_amount)::numeric / 1000000, 2) || ' triệu' as value
FROM orders
UNION ALL
SELECT
  'Total Revenue',
  ROUND(SUM(total_amount)::numeric / 1000000, 2) || ' triệu'
FROM orders
UNION ALL
SELECT
  'Highest Order',
  ROUND(MAX(total_amount)::numeric / 1000000, 2) || ' triệu'
FROM orders
UNION ALL
SELECT
  'Lowest Order',
  ROUND(MIN(total_amount)::numeric / 1000000, 2) || ' triệu'
FROM orders;

-- ===== Review Statistics =====
\echo ''
\echo 'Review Statistics:'
\echo '------------------------------------------'
SELECT
  'Average Rating' as metric,
  ROUND(AVG(rating)::numeric, 2) as value
FROM reviews
UNION ALL
SELECT
  'Total Reviews',
  COUNT(*)::text
FROM reviews
UNION ALL
SELECT
  'Approved %',
  ROUND(100.0 * COUNT(*) FILTER (WHERE is_approved = true) / COUNT(*)::numeric, 1) || '%'
FROM reviews;
