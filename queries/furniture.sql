-- name: ListPopularFurniture :many
SELECT COUNT(*) OVER ()             AS total_records,
       f.id,
       f.name,
       f.price,
       c.name                       AS category,
       COALESCE(AVG(r.rating), 0)   AS average_rating,
       (w.furniture_id IS NOT NULL) AS wishlisted,
       fi.url                       AS image_url
FROM furniture f
         JOIN categories c ON f.category_id = c.id
         LEFT JOIN wishlist w ON w.furniture_id = f.id AND w.identity_id = $1
         LEFT JOIN reviews r ON r.furniture_id = f.id
         LEFT JOIN LATERAL (
            SELECT url
            FROM furniture_images
            WHERE furniture_id = f.id
              AND is_primary = true
            LIMIT 1
    ) fi ON true
WHERE f.stock_quantity > 0
  AND ($2::TEXT IS NULL OR c.name = $2)
  AND ($3::DECIMAL IS NULL OR f.price <= $3)
  AND ($4::INT IS NULL OR r.rating = $4)
GROUP BY f.id, f.name, f.price, c.name, fi.url, w.furniture_id
ORDER BY average_rating DESC, f.price
LIMIT $5 OFFSET $6;

-- name: ListRecommendedFurniture: many
SELECT COUNT(*) OVER ()              as total_records,
       f.id,
       f.name,
       f.price,
       c.name                        AS category,
       COALESCE(AVG(r.rating), 0)    AS average_rating,
       (uw.furniture_id IS NOT NULL) AS wishlisted,
       fi.url                        AS image_url,
       COUNT(w.furniture_id)         AS wishlist_count
FROM furniture f
         JOIN categories c ON f.category_id = c.id
         LEFT JOIN wishlist uw ON uw.furniture_id = f.id AND uw.identity_id = $1
         LEFT JOIN wishlist w ON w.furniture_id = f.id
         LEFT JOIN reviews r ON r.furniture_id = f.id
         LEFT JOIN LATERAL (
            SELECT url
            FROM furniture_images
            WHERE furniture_id = f.id
              AND is_primary = true
            LIMIT 1
        ) fi ON true
WHERE f.stock_quantity > 0
  AND ($2::TEXT IS NULL OR c.name = $2)
  AND ($3::DECIMAL IS NULL OR f.price <= $3)
  AND ($4::INT IS NULL OR r.rating = $4)
GROUP BY f.id, f.name, f.price, c.name, fi.url, uw.furniture_id
ORDER BY wishlist_count DESC, average_rating DESC
LIMIT $5 OFFSET $6;