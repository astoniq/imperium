package objecttype

import "time"

type Spec struct {
	Type      string                  `json:"type"`
	Source    *Source                 `json:"source,omitempty"`
	Relations map[string]RelationRule `json:"relations"`
	CreatedAt time.Time               `json:"createdAt,omitempty"`
}

type Source struct {
	DatabaseType string `json:"dbType" validate:"required"`
}

type CreateSpec struct {
	Type      string                  `json:"type" validate:"required,valid_object_type"`
	Source    *Source                 `json:"source,omitempty"`
	Relations map[string]RelationRule `json:"relations" validate:"required,min=1,dive"`
}

type RelationRule struct {
	InheritIf    string         `json:"inheritIf,omitempty" validate:"required_with=Rules OfType WithRelation,valid_inheritif"`
	Rules        []RelationRule `json:"rules,omitempty"        validate:"required_if_oneof=InheritIf anyOf allOf noneOf,omitempty,min=1,dive"` // Required if InheritIf is "anyOf", "allOf", or "noneOf", empty otherwise
	OfType       string         `json:"ofType,omitempty"       validate:"required_with=WithRelation,valid_relation"`
	WithRelation string         `json:"withRelation,omitempty" validate:"required_with=OfType,valid_relation"`
}
