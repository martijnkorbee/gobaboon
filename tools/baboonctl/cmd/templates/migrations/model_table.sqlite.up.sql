CREATE TABLE IF NOT EXISTS $TABLENAME$ (
	"id" INTEGER,
	"created_at" TEXT NOT NULL DEFAULT 'datetime()',
	"updated_at" TEXT NOT NULL DEFAULT 'datetime()',
	PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE TRIGGER IF NOT EXISTS set_timestamp_$TABLENAME$
	BEFORE UPDATE ON $TABLENAME$
BEGIN
	UPDATE $TABLENAME$ SET updated_at = datetime();
END;
