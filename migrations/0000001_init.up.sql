CREATE TABLE IF NOT EXISTS "user" (
    "id" INTEGER NOT NULL,
    "username" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "firstname" TEXT NOT NULL,
    "lastname" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "failed_attempts" INTEGER NOT NULL DEFAULT 0,
    "security_stamp" TEXT NOT NULL,
    "last_login" TEXT NOT NULL,
    "data" BLOB NOT NULL,
    PRIMARY KEY("id")
);