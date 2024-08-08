INSERT INTO user_sessions (user_id, token, created_at, expires_at, ip_address, user_agent)
VALUES (
    $1, -- 1 user_id
    $2, --705a... Example JWT token
    $3, --CURRENT_TIMESTAMP,  -- Automatically uses the current time
    $4, --CURRENT_TIMESTAMP + interval '7 days',  -- Replace with the desired expiration time
    '', --'192.168.1.1'  -- Replace with the user's IP address
    '', --'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3'  -- Example user agent
);