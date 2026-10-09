-- Fills in category/banner photos that were left blank (or stuck on random
-- picsum.photos placeholders) after the design handoff v2 redesign.
-- Curated, theme-matched stock photos (Unsplash, free-to-use license).
-- Idempotent: only touches rows still on NULL or a picsum.photos placeholder,
-- so re-running this migration never overwrites an image an admin set later.

UPDATE categories SET image_url = 'https://images.unsplash.com/photo-1638517317391-af4c18e4c96a?w=1200&h=800&fit=crop&q=80'
  WHERE id = 'tranh-dong' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE categories SET image_url = 'https://images.unsplash.com/photo-1684871430772-569936b1a0ae?w=1200&h=800&fit=crop&q=80'
  WHERE id = 'tranh-phong-thuy' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE categories SET image_url = 'https://images.unsplash.com/photo-1612704057720-e8f66bade6ca?w=1200&h=800&fit=crop&q=80'
  WHERE id = 'dinh-dong-tho-cung' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE categories SET image_url = 'https://images.unsplash.com/photo-1651085410796-e663860b2b08?w=1200&h=800&fit=crop&q=80'
  WHERE id = 'tuong-dong-trang-tri' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE categories SET image_url = 'https://images.unsplash.com/photo-1541508223081-3f8cd5cc0f4c?w=1200&h=800&fit=crop&q=80'
  WHERE id = 'do-dung-nha-bep' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE categories SET image_url = 'https://images.unsplash.com/photo-1771795639001-a084a0d1045e?w=1200&h=800&fit=crop&q=80'
  WHERE id = 'phu-kien-trang-tri' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE categories SET image_url = 'https://images.unsplash.com/photo-1689259103820-a375e5a30e00?w=1200&h=800&fit=crop&q=80'
  WHERE id = 'tranh-phong-canh-122507' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');

UPDATE banners SET image_url = 'https://plus.unsplash.com/premium_photo-1671749088116-943334aaad48?w=1600&h=900&fit=crop&q=80'
  WHERE id = 'eca9d5a8-def8-421a-a04e-97dcb805a959' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE banners SET image_url = 'https://images.unsplash.com/photo-1511306162219-1c5a469ab86c?w=1600&h=900&fit=crop&q=80'
  WHERE id = '92468648-1b51-4d8e-b7bc-e6929d599a3d' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE banners SET image_url = 'https://images.unsplash.com/photo-1721508490084-1b1de5b230d4?w=1600&h=900&fit=crop&q=80'
  WHERE id = '0a482436-8c5c-4951-93d5-e3a805f7fd70' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE banners SET image_url = 'https://images.unsplash.com/photo-1547555706-54bcf05bbad1?w=1600&h=900&fit=crop&q=80'
  WHERE id = 'b44bf169-44b1-450d-9d70-2d8970710737' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE banners SET image_url = 'https://plus.unsplash.com/premium_photo-1681433412602-db329923c7e3?w=1600&h=900&fit=crop&q=80'
  WHERE id = '0742fbca-b8d2-4499-95c1-af8057d70ad4' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE banners SET image_url = 'https://images.unsplash.com/photo-1769791650175-6858ef4780bb?w=1600&h=900&fit=crop&q=80'
  WHERE id = '0f90acde-d0a2-4cb4-84f3-f4117a1d5232' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE banners SET image_url = 'https://images.unsplash.com/photo-1771073387047-df16b4889412?w=1600&h=900&fit=crop&q=80'
  WHERE id = 'c08398ad-f6c1-4359-99d0-afed8df2641c' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE banners SET image_url = 'https://images.unsplash.com/photo-1612201598945-f66a763965bd?w=1600&h=900&fit=crop&q=80'
  WHERE id = '7c8e813c-dd37-4013-ad12-02a1e7cec0b3' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE banners SET image_url = 'https://images.unsplash.com/photo-1788823524717-a35209426dff?w=1600&h=900&fit=crop&q=80'
  WHERE id = '0a0adfb6-3783-4d93-a097-e31533e36750' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE banners SET image_url = 'https://images.unsplash.com/photo-1680444551079-319c736f4b0a?w=1600&h=900&fit=crop&q=80'
  WHERE id = 'eaa3b64d-7a2b-456d-b553-86cef9c76668' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
UPDATE banners SET image_url = 'https://plus.unsplash.com/premium_photo-1744901573818-59599d9e7070?w=1600&h=900&fit=crop&q=80'
  WHERE id = 'd447e7a8-dac4-49fa-8d9a-0a6084465bfb' AND (image_url IS NULL OR image_url LIKE '%picsum.photos%');
