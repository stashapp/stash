-- 71_custom_fields.up.sql
CREATE TABLE performer_custom_fields (
  performer_id integer NOT NULL,
  field text NOT NULL,
  "value" bytea NOT NULL,
  PRIMARY KEY ("performer_id", "field"),
  foreign key("performer_id") references "performers"("id") on delete CASCADE
);

CREATE INDEX "index_performer_custom_fields_field_value" ON "performer_custom_fields" ("field", "value");
