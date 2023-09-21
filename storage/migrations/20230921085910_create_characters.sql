-- Create "characters" table
CREATE TABLE "public"."characters" ("id" serial NOT NULL, "name" character varying(255) NOT NULL, "damage" integer NOT NULL, "defense" integer NOT NULL, "critical_odds" integer NOT NULL, "critical_loss" integer NOT NULL, "health" integer NOT NULL, "speed" integer NOT NULL, PRIMARY KEY ("id"));
