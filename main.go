package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func readCSV(filename string) ([]map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'          // Set the delimiter to semicolon
	reader.LazyQuotes = true    // Allow lazy quotes
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read all records at once
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file must have at least a header and one data row")
	}

	header := records[0]
	var data []map[string]string

	for i, record := range records[1:] {
		item := make(map[string]string)
		for j, value := range record {
			if j < len(header) {
				item[header[j]] = value
			} else {
				log.Printf("Warning: Line %d has extra field: %s", i+2, value)
			}
		}
		for _, h := range header {
			if _, exists := item[h]; !exists {
				item[h] = ""
			}
		}
		data = append(data, item)
	}

	return data, nil
}

func getRelevantFields(data []map[string]string, fields []string) []map[string]string {

	var outputSlice []map[string]string
	for _, item := range data {
		relevantMap := make(map[string]string)
		for _, field := range fields {
			relevantMap[field] = item[field]
		}
		outputSlice = append(outputSlice, relevantMap)

	}

	return outputSlice

}

func main() {
	filename := "/home/luho/Code/satm/data.csv"

	data, err := readCSV(filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Loaded %d records from CSV\n", len(data))
	relevantfields := []string{"Codice utente", "IdViaggiatore", "TipoViaggiatore"}
	data = getRelevantFields(data, relevantfields)
	if len(data) > 0 {
		fmt.Println("Sample of the first record:")
		jsonData, err := json.MarshalIndent(data[3], "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonData))

		fmt.Println("\nFields in the first record:")
		for key, value := range data[0] {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	// TODO: API endpoints here
}
