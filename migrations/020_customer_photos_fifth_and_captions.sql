-- The homepage gallery is a 5-column grid; with only 4 rows it left a gap
-- and looked left-packed. Adds a 5th photo and shortens all 5 captions to
-- the "{Phòng} · {Thành phố}" style the design calls for. Idempotent: the
-- insert is keyed by a fixed id (no-op if already present), and each
-- caption update only fires while the caption still matches what 019 set.

INSERT INTO customer_photos (id, image_url, caption, sort_order, is_active)
VALUES (
  'cp-van-phong-5',
  'https://plus.unsplash.com/premium_photo-1732721750870-1be1be92299b?w=1000&h=1333&fit=crop&q=80',
  'Góc làm việc · Đà Nẵng',
  5,
  true
)
ON CONFLICT (id) DO NOTHING;

UPDATE customer_photos SET caption = 'Phòng khách · Hà Nội'
  WHERE id = '6b3da688-2025-43a6-a880-fdeb72c4523e'
  AND caption = 'Phòng khách hiện đại, treo tranh nghệ thuật';

UPDATE customer_photos SET caption = 'Phòng thờ · TP.HCM'
  WHERE id = '9fdf675d-2b35-4ad0-a80b-fd119cef48ba'
  AND caption = 'Góc thờ cúng trang nghiêm trong nhà';

UPDATE customer_photos SET caption = 'Sảnh biệt thự · Hải Phòng'
  WHERE id = 'c62fa580-f24c-47e7-90ba-0435b58e3e89'
  AND caption = 'Không gian sang trọng, phong cách cổ điển';

UPDATE customer_photos SET caption = 'Phòng khách · Bắc Ninh'
  WHERE id = '7fd5ba04-c5a4-4c53-91a1-a53a646f9d00'
  AND caption = 'Phòng khách ấm cúng, gần gũi';
