package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"strings"
	"strconv"
	"sort"
)

// Global variables, 
// Lists and Hashmaps are kept global so our output functions can access them
// data.json is first unmarshalled into our lists then data is put into a Hashmap
// We can use the Type, Pokemon, and Move names as the key since in Pokemon these words are all unique
var types Types
var pokemons Pokemons
var moves Moves
var typeMap map[string]Type
var pokemonMap map[string]Pokemon
var moveMap map[string]Move
const dataLocation string = "data.json" // Change this location to the full path of the data.json file if it is not read.

// Structs for our API
// Struct for Pokemon type
type Type struct {
	// Name of the type
	Name string `json:"name"`
	// The effective types, damage multiplize 2x
	EffectiveAgainst []string `json:"effectiveAgainst"`
	// The weak types that against, damage multiplize 0.5x
	WeakAgainst []string `json:"weakAgainst"`
}
// Types struct is the entire list of Types in the JSON file
type Types struct{
	Types []Type `json:"types"`
}
// Struct for Pokemon
type Pokemon struct {
	Number         string   `json:"Number"`
	Name           string   `json:"Name"`
	Classification string   `json:"Classification"`
	TypeI          []string `json:"Type I"`
	TypeII         []string `json:"Type II,omitempty"`
	Weaknesses     []string `json:"Weaknesses"`
	FastAttackS    []string `json:"Fast Attack(s)"`
	Weight         string   `json:"Weight"`
	Height         string   `json:"Height"`
	Candy          struct {
		Name     string `json:"Name"`
		FamilyID int    `json:"FamilyID"`
	} `json:"Candy"`
	NextEvolutionRequirements struct {
		Amount int    `json:"Amount"`
		Family int    `json:"Family"`
		Name   string `json:"Name"`
	} `json:"Next Evolution Requirements,omitempty"`
	NextEvolutions []struct {
		Number string `json:"Number"`
		Name   string `json:"Name"`
	} `json:"Next evolution(s),omitempty"`
	PreviousEvolutions []struct {
		Number string `json:"Number"`
		Name   string `json:"Name"`
	} `json:"Previous evolution(s),omitempty"`
	SpecialAttacks      []string `json:"Special Attack(s)"`
	BaseAttack          int      `json:"BaseAttack"`
	BaseDefense         int      `json:"BaseDefense"`
	BaseStamina         int      `json:"BaseStamina"`
	CaptureRate         float64  `json:"CaptureRate"`
	FleeRate            float64  `json:"FleeRate"`
	BuddyDistanceNeeded int      `json:"BuddyDistanceNeeded"`
}
// Moves struct is the entire list of Pokemon in the JSON file
type Pokemons struct {
	Pokemons []Pokemon `json:"pokemons"`
}

// Move is an attack information a pokemon can have.
type Move struct {
	// The ID of the move
	ID int `json:"id"`
	// Name of the attack
	Name string `json:"name"`
	// Type of attack
	Type string `json:"type"`
	// The damage that enemy will take
	Damage int `json:"damage"`
	// Energy requirement of the attack
	Energy int `json:"energy"`
	// Dps is Damage Per Second
	Dps float64 `json:"dps"`
	// The duration
	Duration int `json:"duration"`
}
// Moves struct is the entire list of Moves in the JSON file
type Moves struct {
	Moves []Move `json:"moves"`
}
// BaseData is a struct for reading data.json
type BaseData struct {
	Types    []Type    `json:"types"`
	Pokemons []Pokemon `json:"pokemons"`
	Moves    []Move    `json:"moves"`
}
// Print functions That can be called by our handler functions. Seperated so they an be reuseable.
func printType( w http.ResponseWriter, pokemonType Type){
	fmt.Fprint( w, "Type Name: " + pokemonType.Name + "\n")
	fmt.Fprint( w, "Effective Against: " + "\n")
	for i:= 0; i < len(pokemonType.EffectiveAgainst); i++{
		fmt.Fprint( w, "	" + pokemonType.EffectiveAgainst[i] + "\n")
	}
	fmt.Fprint( w,  "Weak Against: " + "\n")
	for  i:= 0; i < len(pokemonType.WeakAgainst); i++{
		fmt.Fprint( w, "	" + pokemonType.WeakAgainst[i] + "\n")
	}
}

func printPokemon( w http.ResponseWriter, pokemon Pokemon){
	fmt.Fprint( w, "Poekmon Number: " + pokemon.Number + "\n" + 
					"Name: " + pokemon.Name + "\n" +
					"Classification: " + pokemon.Classification + "\n" +
					"TypeI: " + "\n")
		for i:= 0; i < len(pokemon.TypeI); i++ {
			fmt.Fprint( w, "	" + pokemon.TypeI[i] + "\n")
		}
		fmt.Fprint( w, "TypeII:" + "\n")
		for i:= 0; i < len(pokemon.TypeII); i++ {
			fmt.Fprint( w, "	" + pokemon.TypeII[i] + "\n")
		}
		fmt.Fprint( w, "Weaknesses:" + "\n")
		for i:= 0; i < len(pokemon.Weaknesses); i++ {
			fmt.Fprint( w, "	" + pokemon.Weaknesses[i] + "\n")
		}
		fmt.Fprint( w, "FastAttacks:" + "\n")
		for i:= 0; i < len(pokemon.FastAttackS); i++ {
			fmt.Fprint( w, "	" + pokemon.FastAttackS[i] + "\n")
		}
		fmt.Fprint( w, "Height: " + pokemon.Height + "\n" + 
						"Weight: " + pokemon.Weight + "\n")
		fmt.Fprint( w, "Candy: " + "\n" + 
						"	Name:" + pokemon.Candy.Name + "\n" + 
						"	FamilyID" + strconv.Itoa(pokemon.Candy.FamilyID) + "\n" +
						"NextEvolutionRequirements:" + "\n" +
						"	Amount: " + strconv.Itoa( pokemon.NextEvolutionRequirements.Amount) + "\n" +
						"	Family: " + strconv.Itoa( pokemon.NextEvolutionRequirements.Family) + "\n" +
						"	Name: " + pokemon.NextEvolutionRequirements.Name + "\n" +
						"NextEvolutions:" + "\n")
		for i:= 0; i <len(pokemon.NextEvolutions); i++{
		fmt.Fprint( w, "	Number: " + pokemon.NextEvolutions[i].Number + 
						" ; Name: " + pokemon.NextEvolutions[i].Name + "\n")
		}
		fmt.Fprint( w, "Previous Evolutions: " + "\n")
		for i:= 0; i <len(pokemon.PreviousEvolutions); i++{
		fmt.Fprint( w, "	Number: " + pokemon.PreviousEvolutions[i].Number + 
						" ; Name: " + pokemon.PreviousEvolutions[i].Name + "\n")
		}
		fmt.Fprint( w, "Special Attacks: " + "\n")
		for i:= 0; i <len(pokemon.SpecialAttacks); i++{
		fmt.Fprint( w, "	" + pokemon.SpecialAttacks[i])
		}
		fmt.Fprint( w, "Base Attack: " + strconv.Itoa(pokemon.BaseAttack) + "\n" +
					"Base Defense: " + strconv.Itoa(pokemon.BaseDefense) + "\n" +
					"Base Stamina: " + strconv.Itoa(pokemon.BaseStamina) + "\n" +
					"Capture Rate: " + strconv.FormatFloat(pokemon.CaptureRate, 'f', -1, 64) + "\n" +
					"Flee Rate: " + strconv.FormatFloat(pokemon.FleeRate, 'f', -1, 64) + "\n" +
					"Buddy Distance Needed: " + strconv.Itoa(pokemon.BuddyDistanceNeeded) + "\n")
}

func printMove( w http.ResponseWriter, move Move){
	fmt.Fprint( w, "Move ID: " + strconv.Itoa(move.ID) + "\n" + 
	"Name: " + move.Name + "\n" + 
	"Type: " + move.Type + "\n" + 
	"Damage: " + strconv.Itoa(move.Damage) + "\n" + 
	"Energy:" + strconv.Itoa(move.Energy) + "\n" + 
	"DPS: " + strconv.FormatFloat(move.Dps, 'f', -1, 64) + "\n" + 
	"Duration: " + strconv.Itoa(move.Duration) + "\n") 
}

// Handler functions for HTTP requests
// List Handlers will list pokemons or moves according to what the user wants
func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/list url:", r.URL)
	q:= r.URL.Query() // create a Values map of the Users Query
	p:= make([]Pokemon,0) // Slice of pokemon obkects
	// Look for a type request in query
	// If found go through the entire List of pokemon and print pokemons who have TypeI or TypeII of this type
	s, ok := q["type"]
	if ok {
		for i:= 0; i < len(pokemons.Pokemons); i++{
			for j:=0; j < len(s); j++ {
				if( strings.ToTitle(pokemons.Pokemons[i].TypeI[0]) == strings.ToTitle(s[j]) || strings.ToTitle(pokemons.Pokemons[i].TypeI[0]) == strings.ToTitle(s[j])){
					p = append(p, pokemons.Pokemons[i])
				}
			}

		}
	}
	// Tests sort parameters. 
	// Uses the Built in sort.slice functionality of Go
	sortBy, ok := q["sortby"]
	if ok {
		switch sortBy[0]{
			default: // return an error because it's an invalid value
			case "BaseAttack":
				sort.Slice(p, func(i, j int) bool {
					return p[i].BaseAttack > p[j].BaseAttack
				})
			case "BaseDefense" :
				sort.Slice(p, func(i, j int) bool {
					return p[i].BaseDefense > p[j].BaseDefense
				})
			case "BaseStamina" :
				sort.Slice(p, func(i, j int) bool {
					return p[i].BaseStamina > p[j].BaseStamina
				})
			case "CaptureRate" :
				sort.Slice(p, func(i, j int) bool {
					return p[i].CaptureRate > p[j].CaptureRate
				})
			case "FleeRate" :
				sort.Slice(p, func(i, j int) bool {
					return p[i].FleeRate > p[j].FleeRate
				})
			case "Weight" :
				sort.Slice(p, func(i, j int) bool {
					return p[i].Weight > p[j].Weight
				})
			case "Height" :
				sort.Slice(p, func(i, j int) bool {
					return p[i].Height > p[j].Height
				})
			case "BuddyDistanceNeeded" :
				sort.Slice(p, func(i, j int) bool {
					return p[i].BuddyDistanceNeeded > p[j].BuddyDistanceNeeded
				})
		}
	}
	// Finally Print our data
	for i:= 0; i < len(p); i++{
		printPokemon( w, p[i])
		fmt.Fprint( w, "\n -------------------------------------------------------------------- \n \n")
	}
}

// The API already accesses the user to data in the otherwise function. 
func getHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/get url:", r.URL)

	fmt.Fprint(w, "The Get Handler\n")
}

func getType(w http.ResponseWriter, r *http.Request) {
	log.Println("/list/types url:", r.URL)
	fmt.Fprint(w, "Types\n")
	for i := 0; i < len(types.Types); i++ {
		fmt.Fprint(w, "Type name: " + types.Types[i].Name + "\n")
	}
}
// Otherwise function for single word or to return a non indexed word
func otherwise(w http.ResponseWriter, r *http.Request) {
	// Data holders for our mapped JSON data
	var pokemonType Type
	var pokemon Pokemon
	var move Move
	// String to read what the user wants
	str := strings.TrimLeft( r.URL.String(), "/")
	// Change all instances of "%20" to a blank space " " so the word matches what we have in our data
	// Spaces in URL are defaulted to "%20"
	str = strings.Replace(str, "%20", " ", -1) 
	// Boolean to indicate if we have an object the user wants or not
	found := false
	// Finds if the word the user is looking for is a pokemon type then calls the pokemon type output function.
	pokemonType, ok:= typeMap[strings.ToTitle(str)]
	if( ok){
		printType( w, pokemonType)
		found = ok // We have what the user is looking for, set to true. 
	}
	// Finds if the word the user is looking for is a pokemon then calls the pokemon output function.
	pokemon, ok = pokemonMap[strings.ToTitle(str)]
	if(ok){
		printPokemon( w, pokemon)
		found = ok // We have what the user is looking for, set to true. 
	}
	// Finds if the word the user is looking for is a move then calls the move output function.
	move, ok = moveMap[strings.ToTitle(str)]
	if( ok){
		printMove( w, move)
		found = ok // We have what the user is looking for, set to true. 
	}
	// Enters if we could not find what the user is looking for
	if !found{
		fmt.Fprint( w, "We could not find what you are looking for")
	}
}

func main() {

	// Open our jsonFile
	jsonFile, err := os.Open( dataLocation)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully Opened data.json")
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Types array

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &types)
	json.Unmarshal(byteValue, &pokemons)
	json.Unmarshal(byteValue, &moves)

	// Put our data into a hashmap
	// Intialize the map then iterate through our arrays
	typeMap = make(map[string]Type)
	pokemonMap = make(map[string]Pokemon)
	moveMap = make(map[string]Move)
	// Put our data types into their respective Hashmaps.
	for i := 0; i < len(types.Types); i++ {
		typeMap[ strings.ToTitle(types.Types[i].Name)] = types.Types[i]
	}

	for i:= 0; i < len(moves.Moves); i++ {
		moveMap[ strings.ToTitle(moves.Moves[i].Name)] = moves.Moves[i]
	}

	for i:= 0; i < len(pokemons.Pokemons); i++ {
		pokemonMap[ strings.ToTitle(pokemons.Pokemons[i].Name)] = pokemons.Pokemons[i]
	}
	// Handle Functions are ready to be called
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/list/types", getType)
	http.HandleFunc("/", otherwise)
	log.Println("starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
