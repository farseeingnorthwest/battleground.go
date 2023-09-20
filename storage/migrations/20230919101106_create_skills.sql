-- Create "skills" table
CREATE TABLE "public"."skills" ("id" serial NOT NULL, "name" character varying(255) NOT NULL, "reactor" jsonb NOT NULL, PRIMARY KEY ("id"));
