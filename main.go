package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

func idViaggiatoreBuilder(data []map[string]string) {

	relevantFields := []string{"Codice NUTS PAB", "Codice azienda", "Numero viaggiatore"}
	subsetData := getRelevantFields(data, relevantFields)

	for i, a := range data {
		for _, value := range subsetData[i] {
			a["IdViaggiatore"] += value
		}
		a["IdViaggiatore"] += "LUHO"
	}

}

func tipoViaggiatoreBuilder(data []map[string]string) {

	relevantFields := []string{"Codice ISTAT regione", "Sigla vecchie targe auto", "Codice ISTAT PAB", "Codice ISTAT Comune domicilio",
		"Codice ISTAT CAP domicilio", "Universo", "Genere2", "Condizione occupazione", "Professione", "Età3", "Diversa abilità4", "Altre limitazioni_NoDisabilita", "Altre limitazioni_SiDisabilita"}
	subsetData := getRelevantFields(data, relevantFields)

	for i, a := range data {
		for _, value := range subsetData[i] {
			a["TipoViaggiatore"] += value
		}
		a["TipoViaggiatore"] += "LUHO"
	}

}

func getUsers(data []map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, data)
	}
}

func getUserById(data []map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("CodiceUtente")
		fmt.Println("WHERE  ", id)

		for _, a := range data {
			if id == a["Codice utente"] {
				c.IndentedJSON(http.StatusOK, a)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
	}
}

func main() {
	filename := "/home/luho/Code/satm/data.csv"

	data, err := readCSV(filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Loaded %d records from CSV\n", len(data))
	relevantfields := []string{"Codice utente", "IdViaggiatore", "TipoViaggiatore"}
	idViaggiatoreBuilder(data)
	tipoViaggiatoreBuilder(data)
	data = getRelevantFields(data, relevantfields)

	// TODO: API endpoints here
	router := gin.Default()
	router.GET("/users", getUsers(data))
	router.GET("/users/:CodiceUtente", getUserById(data))

	router.Run("localhost:8080")
}
