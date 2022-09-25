//go:build tools
// +build tools

package main

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
)

var names map[string]*naming

type performerNamingConfig struct {
	Male    string `yaml:"male"`
	Female  string `yaml:"female"`
	Surname string `yaml:"surname"`
}

type namingConfig struct {
	Scenes     string                `yaml:"scenes"`
	Performers performerNamingConfig `yaml:"performers"`
	Galleries  string                `yaml:"galleries"`
	Studios    string                `yaml:"studios"`
	Images     string                `yaml:"images"`
	Tags       string                `yaml:"tags"`
}

type naming struct {
	names []string
}

func (n naming) generateName(words int) string {
	var ret []string
	for i := 0; i < words; i++ {
		w := rand.Intn(len(n.names))
		ret = append(ret, n.names[w])
	}

	return strings.Join(ret, " ")
}

func createNaming(fn string) (*naming, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := &naming{}
	s := bufio.NewScanner(file)
	for s.Scan() {
		ret.names = append(ret.names, s.Text())
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func initNaming(c config) {
	names = make(map[string]*naming)
	load := func(v string) {
		if names[v] == nil {
			var err error
			names[v], err = createNaming(v)
			if err != nil {
				panic(err)
			}
		}
	}

	n := c.Naming
	load(n.Galleries)
	load(n.Images)
	load(n.Scenes)
	load(n.Studios)
	load(n.Tags)
	load(n.Performers.Female)
	load(n.Performers.Male)
	load(n.Performers.Surname)
}

func generatePerformerName() string {
	female := rand.Intn(4) > 0
	wordRand := rand.Intn(100)
	givenNames := 1
	surnames := 1
	if wordRand < 3 {
		givenNames = 2
	} else if wordRand < 26 {
		surnames = 0
	}

	fn := c.Naming.Performers.Female
	if !female {
		fn = c.Naming.Performers.Male
	}

	name := names[fn].generateName(givenNames)
	if surnames > 0 {
		name += " " + names[c.Naming.Performers.Surname].generateName(1)
	}

	return name
}
