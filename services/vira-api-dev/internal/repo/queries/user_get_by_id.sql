SELECT user_id, city, joined_at
FROM user_profiles
WHERE user_id = $1;