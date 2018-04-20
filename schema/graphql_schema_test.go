package schema_test

import (
	"fmt"
	"testing"

	"github.com/suppayami/goql/schema"
)

func TestGraphqlFieldStringer(t *testing.T) {
	gqlField := schema.GraphqlField{
		Name:       "hero",
		Type:       schema.ObjectType,
		ObjectType: "Character",
		Nullable:   true,
		Arguments: []schema.GraphqlArgument{
			schema.GraphqlArgument{
				Name:         "episode",
				Type:         schema.ObjectType,
				ObjectType:   "Episode",
				Nullable:     false,
				DefaultValue: "NEWHOPE",
			},
		},
	}

	expected := "hero(episode: Episode = NEWHOPE): Character"

	if gqlField.String() != expected {
		t.Fatal(fmt.Sprintf("Expected: \n%s\nGot:\n%s\n", expected, gqlField.String()))
	}
}
