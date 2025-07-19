package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"time"
)

// Shorturl holds the schema definition for the Shorturl entity.
type Shorturl struct {
	ent.Schema
}

// Fields of the Shorturl.
func (Shorturl) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			Annotations(
				entsql.Annotation{
					//Incremental: true, // 标记为自增
				},
			).
			Comment("自增主键"),
		field.String("short_code").
			Unique().
			Immutable().
			MaxLen(20).
			Comment("短码(6-8位字符)"),

		field.String("long_url").
			NotEmpty().
			MaxLen(2048).
			SchemaType(map[string]string{
				"mysql": "varchar(2048) COLLATE utf8mb4_bin",
			}).
			Comment("原始URL"),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("创建时间"),

		field.Time("expire_at").
			Optional().
			Nillable().
			Comment("过期时间(null表示永久有效)"),

		field.Bool("is_deleted").
			Default(false).
			Comment("软删除标记"),

		field.Int("access_count").
			Default(0).
			Comment("访问次数"),

		field.String("creator_ip").
			Optional().
			MaxLen(45).
			Comment("创建者IP"),

		field.String("creator_id").
			Optional().
			MaxLen(64).
			Comment("创建者用户ID"),
	}
}

// Mixin 混入基础字段
func (Shorturl) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Annotations 表级注释
func (Shorturl) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
			Options:   "COMMENT='短链映射表'",
		},
	}
}

// Edges of the Shorturl.
func (Shorturl) Edges() []ent.Edge {
	return nil
}
