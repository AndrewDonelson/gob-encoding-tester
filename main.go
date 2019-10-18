package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

// Persist defines the functions needed to store and retreive data
type Persist interface {
	Init() error
	Load()
	Save()
}

// Encoding allows encoding & decoding using gob
type Encoding interface {
	Encode()
	Decode()
}

// Animals defines all interfaces and functions needed to implement an Animal object
type Animals interface {
	//Encoding
	Persist
	Birth(name, kind string)
	Speak()
	Fetch(what, where string)
}

// PersistBoltDB provides BoldDB storage
type PersistBoltDB struct {
	db     *bolt.DB
	bucket *bolt.Bucket
}

func (p *PersistBoltDB) initBucket(name string) (err error) {

	p.db.Update(func(tx *bolt.Tx) error {
		p.bucket, err = tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			err = fmt.Errorf("create bucket %s: %s", name, err)
		}
		return err
	})
	return
}

// Init initializes & opens a BolDB instance
func (p *PersistBoltDB) Init(kind string) (err error) {

	// initialize database
	p.db, err = bolt.Open("animal.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("Error with DB initialization. -> %v", err)
	}
	defer p.db.Close()

	// initialize bucket for this Animal Type (Cat & Dog uses different buckets)
	p.initBucket(kind)
	return nil
}

// Load gets a object from BoltDB
func (p *PersistBoltDB) Load(key string) {
	p.db.View(func(tx *bolt.Tx) error {
		//b := tx.Bucket([]byte("MyBucket"))
		v := p.bucket.Get([]byte(key))
		fmt.Printf("Loaded: %s\n", v)
		return nil
	})
}

// Save puts a object into BoltDB
func (p *PersistBoltDB) Save(key string, value interface{}) {
	b, ok := value.(*[]byte)
	if ok {
		p.db.Update(func(tx *bolt.Tx) error {
			//b := tx.Bucket([]byte("MyBucket"))
			err := p.bucket.Put([]byte(key), *b)
			return err
		})
	}
}

// ----- Animal ---------------------------------------//

// Animal defines the common features for all animals.
// Must implement all Action methods at a minimum
// Do not create new instances of Animal directly. Instead derive a new Type,
// say Fish and add Animal there. See Cat & Doc types as example
type Animal struct {
	Kind        string
	Name        string
	FavoriteToy string
}

// setKind is a helper function. Use in call from New???() with `reflect.TypeOf(m).String()` as the parameter for typeName
func (a *Animal) setKind(typeName string) {
	// Get the short name (Animal instead of main.Animal) for this KIND of animal. ie. (Cat, Dog)
	s := strings.Split(typeName, ".")
	a.Kind = s[len(s)-1]
}

// Birth is used to create a new Animal
// name is the name of the animal ("Spot, Kitty")
// kind is the Type for the Animal (Cat, Dog) - Use `reflect.TypeOf(m).String()` as the parameter for kind
func (a *Animal) Birth(name, kind string) {
	a.Name = name
	a.FavoriteToy = ""
	a.setKind(kind)
	fmt.Printf("A %s named %s was born!\n", a.Kind, a.Name)
}

// Speak allows the animal to speak
func (a *Animal) Speak(what string) {
	fmt.Printf("%s says %s\n", a.Name, what)
}

// Fetch is implemented in Derrived Animal Type (cat, Dog, etc)
func (a *Animal) Fetch(what, where string) {
	fmt.Println("")
}

// Init initializes & opens a Persist instance
func (a *Animal) Init() error {
	return nil
}

// Load gets an animal from Persist object
func (a *Animal) Load() {

}

// Save puts an animal into  Persist object
func (a *Animal) Save() {

}

// ----- Cat ---------------------------------------//

// Cat defines Animal of Type Cat
type Cat struct {
	Animal
}

// Speak allows the Cat to speak
func (c *Cat) Speak() {
	c.Animal.Speak("Meow")
}

// Fetch allows the Cat to fetch something from somewhere
func (c *Cat) Fetch(what, where string) {
	fmt.Printf("Yeah...Right, %s's don't fetch!\n", c.Kind)
}

// Init initializes & opens a Persist instance
func (c *Cat) Init() error {
	return nil
}

// Load gets an Cat from a Persist object
func (c *Cat) Load() {

}

// Save stores a Cat to a Persist object
func (c *Cat) Save() {

}

// NewCat creates an new Cat object and returns a reference to it
func NewCat(name string) *Cat {
	m := new(Cat)
	m.Birth(name, reflect.TypeOf(m).String())
	return m
}

// ----- Dog ---------------------------------------//

// Dog defines Animal of Type Dog
type Dog struct {
	Animal
}

// Speak allows the Dog to speak
func (d *Dog) Speak() {
	d.Animal.Speak("Woof")
}

// Fetch allows the Dog to fetch something from somewhere
func (d *Dog) Fetch(what, where string) {
	fmt.Printf("Your favorite %s %s, Happily fetches %s from the %s\n", d.Kind, d.Name, what, where)
}

// NewDog creates an new Dog object and returns a reference to it
func NewDog(name string) *Dog {
	m := new(Dog)
	m.Birth(name, reflect.TypeOf(m).String())
	return m
}

func initialize() (enc *gob.Encoder, dec *gob.Decoder) {
	// This type must match exactly what youre going to be using,
	// down to whether or not its a pointer
	gob.Register(&Cat{})
	gob.Register(&Dog{})

	network := new(bytes.Buffer)
	enc = gob.NewEncoder(network)
	dec = gob.NewDecoder(network)

	return
}

func main() {
	enc, dec := initialize()

	var inter Animals
	inter = NewCat("Garfield")
	// Note: pointer to the interface
	err := enc.Encode(&inter)
	if err != nil {
		panic(err)
	}

	inter = NewDog("Snoopy")
	// Note: pointer to the interface
	err = enc.Encode(&inter)
	if err != nil {
		panic(err)
	}

	// Now lets get them back out

	var get Animals
	err = dec.Decode(&get)
	if err != nil {
		panic(err)
	}
	// Should Meow
	get.Speak()
	get.Fetch("mouse", "closet")

	err = dec.Decode(&get)
	if err != nil {
		panic(err)
	}
	// Should Woof
	get.Speak()
	get.Fetch("ball", "yard")

}
