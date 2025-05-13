SELECT *, row_to_json(m1.*) as logo, row_to_json(m2.*) as cover
FROM organizations o
JOIN media m1 ON m1.id=o.image
JOIN media m2 ON m2.id=o.cover_image
WHERE id IN(?)