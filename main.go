package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Persist interface {
	Init()
	Load()
	Save()
}

type Actions interface {
	Birth(name string)
	Speak()
}

type PersistBoltDB int

func (p *PersistBoltDB) Init() error { return nil }
func (p *PersistBoltDB) Load()       {}
func (p *PersistBoltDB) Save()       {}

type Animal struct {
	Name string
}

func (a *Animal) Birth(name string) {
	a.Name = name
	fmt.Printf("%s was born!\n", a.Name)
}
func (a *Animal) Speak(what string) {
	fmt.Printf("%s says %s\n", a.Name, what)
}
func (a *Animal) Load() {

}
func (a *Animal) Save() {

}

type Cat struct {
	Animal
}

type Dog struct {
	Animal
}

func (c *Cat) Speak() {
	c.Animal.Speak("Meow")
}
func (c *Cat) Load() {

}
func (c *Cat) Save() {

}

func (d *Dog) Speak() {
	d.Animal.Speak("Woof")
}

func init() {
	// This type must match exactly what youre going to be using,
	// down to whether or not its a pointer
	gob.Register(&Cat{})
	gob.Register(&Dog{})
}

func main() {
	network := new(bytes.Buffer)
	enc := gob.NewEncoder(network)

	var inter Actions
	inter = new(Cat)
	inter.Birth("Garfield")

	// Note: pointer to the interface
	err := enc.Encode(&inter)
	if err != nil {
		panic(err)
	}

	inter = new(Dog)
	inter.Birth("Snoopy")
	err = enc.Encode(&inter)
	if err != nil {
		panic(err)
	}

	// Now lets get them back out
	dec := gob.NewDecoder(network)

	var get Actions
	err = dec.Decode(&get)
	if err != nil {
		panic(err)
	}

	// Should meow
	get.Speak()

	err = dec.Decode(&get)
	if err != nil {
		panic(err)
	}

	// Should woof
	get.Speak()

}
