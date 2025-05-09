SELECT oc.*
FROM users u
JOIN oauth_connects oc ON oc.identity_id=u.id
WHERE u.email=$1 AND oc.provider=$2