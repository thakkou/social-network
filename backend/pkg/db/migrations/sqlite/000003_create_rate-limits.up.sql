--Rate Limits
CREATE TABLE IF NOT EXISTS rate_limits (
    ip TEXT NOT NULL,
    route TEXT NOT NULL,
    last_request DATETIME NOT NULL,
    PRIMARY KEY (ip, route)
);