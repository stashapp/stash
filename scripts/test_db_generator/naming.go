package main

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
)

var names map[string]*naming

type performerNamingConfig struct {
	male    string `yaml:"male"`
	female  string `yaml:"female"`
	surname string `yaml:"surname"`
}

type namingConfig struct {
	scenes     string                `yaml:"scenes"`
	performers performerNamingConfig `yaml:"performers"`
	galleries  string                `yaml:"galleries"`
	studios    string                `yaml:"studios"`
	images     string                `yaml:"images"`
	tags       string                `yaml:"tags"`
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
	file, err := os.Open("config.yml")
	defer file.Close()
	if err != nil {
		return nil, err
	}

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

	n := c.naming
	load(n.galleries)
	load(n.images)
	load(n.scenes)
	load(n.studios)
	load(n.tags)
	load(n.performers.female)
	load(n.performers.male)
	load(n.performers.surname)
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

	fn := c.naming.performers.female
	if !female {
		fn = c.naming.performers.male
	}

	name := names[fn].generateName(givenNames)
	if surnames > 0 {
		name += " " + names[c.naming.performers.surname].generateName(1)
	}

	return name
}
