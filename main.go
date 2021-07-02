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

func main() {
	fmt.Println("Fetching state....")

	// Get state directory
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dirname)

	stateDirectory := dirname + "/.stopit"
	files, err := os.ReadDir(stateDirectory)
	if err != nil {
		fmt.Println("No directory at ~/.stopit. Creating one now.")
		if createDirErr := os.Mkdir(stateDirectory, 0770); createDirErr != nil {
			log.Fatal(createDirErr)
		}
	}
	fmt.Println("Files:")
	fmt.Println(files)

	stateFile := stateDirectory + "/stopItState.txt"

	// Get byte state data if it exists.
	state, err := os.ReadFile(stateFile)
	if os.IsNotExist(err) {
		fmt.Println("Stop it state file doesn't exist, creating one....")
		state = make([]byte, 0)
	}

	// Create starting counter list from byte state data.
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

	// Process command line arguments
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createName := createCmd.String("name", "", "name of counter to create")

	resetCmd := flag.NewFlagSet("reset", flag.ExitOnError)
	resetName := resetCmd.String("name", "", "name of counter to reset")

	if len(os.Args) < 2 {
		fmt.Println("Expected some command. Type -h for help text. Exiting...")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		fmt.Println("Listing counters...")
		for _, c := range counters {
			timeDuration := time.Now().Sub(c.StartTime)
			fmt.Printf("- Counter: %s, Duration: %s \n", c.Name, timeDuration)
		}
	case "create":
		createCmd.Parse(os.Args[2:])
		fmt.Printf("Creating counter with name %s", *createName)
		newCounter := Counter{
			Name:      *createName,
			StartTime: time.Now(),
		}
		counters = append(counters, newCounter)
	case "reset":
		resetCmd.Parse(os.Args[2:])
		fmt.Printf("Reseting counter with name %s", *resetName)
	default:
		fmt.Printf("No such command - %s", command)
	}

	// Serialize and write back to the file
	sFile, err := os.Create(stateFile)
	check(err)
	defer sFile.Close()

	w := bufio.NewWriter(sFile)
	enc := json.NewEncoder(w)
	for _, c := range counters {
		err = enc.Encode(c)
	}

	w.Flush()
}
