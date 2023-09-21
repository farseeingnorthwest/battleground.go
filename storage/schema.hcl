table "characters" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "name" {
    null = false
    type = character_varying(255)
  }
  column "damage" {
    null = false
    type = integer
  }
  column "defense" {
    null = false
    type = integer
  }
  column "critical_odds" {
    null = false
    type = integer
  }
  column "critical_loss" {
    null = false
    type = integer
  }
  column "health" {
    null = false
    type = integer
  }
  column "speed" {
    null = false
    type = integer
  }
  primary_key {
    columns = [column.id]
  }
}
table "skills" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "name" {
    null = false
    type = character_varying(255)
  }
  column "reactor" {
    null = false
    type = jsonb
  }
  primary_key {
    columns = [column.id]
  }
}
schema "public" {
  comment = "standard public schema"
}
