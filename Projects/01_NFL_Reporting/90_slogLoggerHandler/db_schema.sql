CREATE TABLE IF NOT EXISTS logs (
	logTimestamp TIMESTAMP NOT NULL,
	logLevel varchar(50) NOT NULL,
	logMessage varchar(200) NOT NULL,
	fullLog JSONB NOT NULL
);