-- The 4 seeded customer_photos rows were unrelated stock photos (seagulls,
-- a lighthouse, a coastal field, a silhouette portrait) with captions
-- literally labeled "demo"/"mẫu" (sample). Swaps in home-interior photos
-- that actually fit the "Tranh trong nhà khách hàng" section — living
-- rooms with framed wall art, and a household altar (matches the Đỉnh
-- Đồng Thờ Cúng category). Idempotent: only touches rows still on the
-- known placeholder URLs, so it never overwrites real photos an admin
-- has since uploaded.

UPDATE customer_photos SET
  image_url = 'https://images.unsplash.com/photo-1691036561573-4b76998b60de?w=1000&h=1333&fit=crop&q=80',
  caption = 'Phòng khách hiện đại, treo tranh nghệ thuật'
  WHERE id = '6b3da688-2025-43a6-a880-fdeb72c4523e'
  AND image_url LIKE '%x2ndfhoxhisfyasgutqj%';

UPDATE customer_photos SET
  image_url = 'https://images.unsplash.com/photo-1769647097493-b2b891eb21bb?w=1000&h=1333&fit=crop&q=80',
  caption = 'Góc thờ cúng trang nghiêm trong nhà'
  WHERE id = '9fdf675d-2b35-4ad0-a80b-fd119cef48ba'
  AND image_url LIKE '%yepnxkwdmpqvnwj7fq9x%';

UPDATE customer_photos SET
  image_url = 'https://images.unsplash.com/photo-1759774313632-854207c22ec1?w=1000&h=1333&fit=crop&q=80',
  caption = 'Không gian sang trọng, phong cách cổ điển'
  WHERE id = 'c62fa580-f24c-47e7-90ba-0435b58e3e89'
  AND image_url LIKE '%pbx1ut4bxzs1dh1wkczz%';

UPDATE customer_photos SET
  image_url = 'https://images.unsplash.com/photo-1757792859308-b8d2115a0010?w=1000&h=1333&fit=crop&q=80',
  caption = 'Phòng khách ấm cúng, gần gũi'
  WHERE id = '7fd5ba04-c5a4-4c53-91a1-a53a646f9d00'
  AND image_url LIKE '%cjwdurbkspu8lpuronkn%';
