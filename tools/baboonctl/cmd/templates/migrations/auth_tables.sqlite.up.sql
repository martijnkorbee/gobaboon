drop table if exists users;

CREATE TABLE "users" (
	"id" INTEGER,
	"first_name" TEXT NOT NULL,
	"last_name" TEXT NOT NULL,
	"user_active" INTEGER NOT NULL DEFAULT 0,
	"email" TEXT NOT NULL UNIQUE,
	"password" TEXT NOT NULL,
	"created_at" TEXT NOT NULL DEFAULT 'datetime()',
	"updated_at" TEXT NOT NULL DEFAULT 'datetime()',
	PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE TRIGGER IF NOT EXISTS set_timestamp_usr
	BEFORE UPDATE ON users
BEGIN
	UPDATE users SET updated_at = datetime();
END;

drop table if exists remember_tokens;

CREATE TABLE "remember_tokens" (
	"id" INTEGER,
	"user_id" INTEGER NOT NULL,
	"remember_token" TEXT NOT NULL,
	"created_at" TEXT NOT NULL DEFAULT 'datetime()',
	"updated_at" TEXT NOT NULL DEFAULT 'datetime()',
    "expiry" TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TRIGGER IF NOT EXISTS set_timestamp_rmtk
	BEFORE UPDATE ON users
BEGIN
	UPDATE users SET updated_at = datetime();
END;

drop table if exists tokens;

CREATE TABLE "tokens" (
	"id"	INTEGER,
	"user_id"	INTEGER NOT NULL,
	"first_name"	TEXT NOT NULL,
	"email"	TEXT NOT NULL,
	"token"	TEXT NOT NULL,
	"token_hash"	BLOB NOT NULL,
	"created_at"	TEXT NOT NULL DEFAULT 'datetime()',
	"updated_at"	TEXT NOT NULL DEFAULT 'datetime()',
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TRIGGER IF NOT EXISTS set_timestamp_tk
	BEFORE UPDATE ON tokens
BEGIN
	UPDATE tokens SET updated_at = datetime();
END;
