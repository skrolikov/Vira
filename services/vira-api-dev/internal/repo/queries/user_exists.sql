SELECT EXISTS(SELECT 1 FROM user_profiles WHERE user_id = $1);
