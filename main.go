package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type Counter struct {
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Creates a state file directory if one does not exist. Returns the path of the stateFile.
func getOrCreateStateFileLocation() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	stateDirectory := dirname + "/.stopit"
	_, err = os.ReadDir(stateDirectory)
	if err != nil {
		fmt.Println("No directory at $HOME/.stopit Creating one now...")
		if createDirErr := os.Mkdir(stateDirectory, 0770); createDirErr != nil {
			panic(createDirErr)
		}
	}

	stateFile := stateDirectory + "/stopItState.txt"
	return stateFile
}

// Fetch the current state from the statefile or return an emtpty byte slice if there is no current state.
func getCurrentCounters(stateFilePath string) []Counter {
	state, err := os.ReadFile(stateFilePath)
	if os.IsNotExist(err) {
		fmt.Println("Stop it state file doesn't exist, creating one....")
		state = make([]byte, 0)
	}

	counters := make([]Counter, 0)
	dec := json.NewDecoder(bytes.NewReader(state))
	var counter Counter
	for dec.More() {
		err := dec.Decode(&counter)
		if err != nil {
			panic(err)
		}
		counters = append(counters, counter)
	}
	return counters
}

// List current counters
func listCounters(counters []Counter) {
	fmt.Println("Listing counters...")
	for _, c := range counters {
		timeDuration := time.Now().Sub(c.StartTime)
		fmt.Printf("- Counter: %s, Duration: %s \n", c.Name, timeDuration)
	}
}

// Create a new counter and return a slice with that counter included
func createNewCounter(counters []Counter, newCounterName string) []Counter {
	fmt.Printf("Creating counter with name %s", newCounterName)
	if newCounterName == "" {
		fmt.Println("Cannot add counter with blank name")
		return counters
	}
	newCounter := Counter{
		Name:      newCounterName,
		StartTime: time.Now(),
	}
	return append(counters, newCounter)
}

// Create a new counter list without any counter that matches the provided name.
func deleteCounter(counters []Counter, nameToDelete string) []Counter {
	fmt.Printf("Deleting counter with name %s", *&nameToDelete)
	cleanedCounters := make([]Counter, 0, len(counters))
	for _, c := range counters {
		if c.Name != nameToDelete {
			cleanedCounters = append(cleanedCounters, c)
		}
	}
	return cleanedCounters
}

// serialize and write back to file
func writeCurrentCounters(stateFilePath string, counters []Counter) {
	sFile, err := os.Create(stateFilePath)
	check(err)
	defer sFile.Close()

	w := bufio.NewWriter(sFile)
	enc := json.NewEncoder(w)
	for _, c := range counters {
		err = enc.Encode(c)
	}

	w.Flush()
}

func main() {
	fmt.Println("Fetching state....")

	stateFilePath := getOrCreateStateFileLocation()
	counters := getCurrentCounters(stateFilePath)

	// Process command line arguments
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createName := createCmd.String("name", "", "name of counter to create")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteName := deleteCmd.String("name", "", "name of counter to delete")

	if len(os.Args) < 2 {
		fmt.Println("Expected some command. Type -h for help text. Exiting...")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		listCounters(counters)
	case "create":
		createCmd.Parse(os.Args[2:])
		counters = createNewCounter(counters, *createName)
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		counters = deleteCounter(counters, *deleteName)
	default:
		fmt.Printf("No such command - %s", command)
	}

	writeCurrentCounters(stateFilePath, counters)
}
