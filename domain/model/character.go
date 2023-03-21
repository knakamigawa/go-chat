package model

import (
	"fmt"
	"go-chat-ai-server/domain/model/character"
)

type Characters []Character

func (c Characters) Names() []string {
	names := make([]string, len(c))
	fmt.Printf("NAMES %s\n", c)
	for i, v := range c {
		fmt.Printf("NAME %s\n", v.name.String())
		names[i] = v.name.String()
	}
	return names
}

type Character struct {
	name character.Name
	bio  character.Bio
}

func MakeCharacter(name character.Name, bio character.Bio) Character {
	return Character{
		name: name,
		bio:  bio,
	}
}

func (c Character) Name() character.Name {
	return c.name
}

func (c Character) Bio() character.Bio {
	return c.bio
}
