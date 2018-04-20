package schema_test

import (
	"fmt"
	"testing"

	"github.com/suppayami/goql/schema"
)

var (
	gqlHeroField = schema.GraphqlField{
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

	gqlHumanType = schema.GraphqlObjectType{
		Name: "Human",
		Fields: []schema.GraphqlField{
			schema.GraphqlField{
				Name:     "id",
				Type:     schema.ScalarID,
				Nullable: false,
			},

			schema.GraphqlField{
				Name:     "name",
				Type:     schema.ScalarString,
				Nullable: false,
			},
		},
	}
)

func TestGraphqlFieldStringer(t *testing.T) {
	expected := "hero(episode: Episode = NEWHOPE): Character"

	if gqlHeroField.String() != expected {
		t.Fatal(fmt.Sprintf("Expected: \n%s\nGot:\n%s\n", expected, gqlHeroField.String()))
	}
}

func TestGraphqlObjectTypeStringer(t *testing.T) {
	expected := "type Human {\n\tid: ID!\n\tname: String!\n}"

	if gqlHumanType.String() != expected {
		t.Fatal(fmt.Sprintf("Expected: \n%s\nGot:\n%s\n", expected, gqlHumanType.String()))
	}
}
