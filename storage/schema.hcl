table "character_skills" {
  schema = schema.public
  column "character_id" {
    null = false
    type = integer
  }
  column "slot" {
    null = false
    type = smallint
  }
  column "skill_id" {
    null = false
    type = integer
  }
  primary_key {
    columns = [column.character_id, column.slot, column.skill_id]
  }
  foreign_key "characters_skills_character_id_fkey" {
    columns     = [column.character_id]
    ref_columns = [table.characters.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "characters_skills_skill_id_fkey" {
    columns     = [column.skill_id]
    ref_columns = [table.skills.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}
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
