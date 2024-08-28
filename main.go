package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type AIModelConnector struct {
	Client *http.Client
}

type MistralAIModelConnector struct {
	Client *http.Client
	Token  string
}

type Inputs struct {
	Table map[string][]string `json:"table"`
	Query string              `json:"query"`
}

type Response struct {
	Answer      string   `json:"answer"`
	Coordinates [][]int  `json:"coordinates"`
	Cells       []string `json:"cells"`
	Aggregator  string   `json:"aggregator"`
}

func CsvToSlice(data string) (map[string][]string, error) {
	reader := csv.NewReader(strings.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	if len(records) == 0 {
		return result, nil
	}

	headers := records[0]

	for _, header := range headers {
		result[header] = make([]string, len(records)-1)
	}

	for i, row := range records[1:] {
		for j, cell := range row {
			result[headers[j]][i] = cell
		}
	}

	return result, nil
}

func (c *MistralAIModelConnector) GenerateRecommendation(question string) (string, error) {
	url := "https://api-inference.huggingface.co/models/mistralai/Mixtral-8x7B-Instruct-v0.1"

	payload := map[string]interface{}{
		"inputs": question,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rawBody bytes.Buffer
	rawBody.ReadFrom(resp.Body)

	var resultArray []map[string]interface{}
	if err := json.NewDecoder(&rawBody).Decode(&resultArray); err == nil {
		if len(resultArray) > 0 {
			if answer, ok := resultArray[0]["generated_text"].(string); ok {
				return answer, nil
			}
			if answer, ok := resultArray[0]["answer"].(string); ok {
				return answer, nil
			}
		}
	}

	rawBody.Reset()
	rawBody.ReadFrom(resp.Body)

	var resultObject map[string]interface{}
	if err := json.NewDecoder(&rawBody).Decode(&resultObject); err == nil {
		if answer, ok := resultObject["generated_text"].(string); ok {
			return answer, nil
		}
		if answer, ok := resultObject["answer"].(string); ok {
			return answer, nil
		}
	}

	return "", fmt.Errorf("unexpected response format: %s", rawBody.String())
}

func (c *AIModelConnector) ConnectAIModel(payload Inputs, token string) (Response, error) {
	url := "https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq"

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return Response{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Response{}, errors.New("unexpected status code")
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Response{}, err
	}

	aggregator := ""
	if val, ok := result["aggregator"]; ok {
		aggregator = val.(string)
	}

	var coordinates [][]int
	coordinatesRaw := result["coordinates"].([]interface{})
	for _, c := range coordinatesRaw {
		coords := c.([]interface{})
		var pair []int
		for _, coord := range coords {
			pair = append(pair, int(coord.(float64)))
		}
		coordinates = append(coordinates, pair)
	}

	var cells []string
	cellsRaw := result["cells"].([]interface{})
	for _, c := range cellsRaw {
		cells = append(cells, c.(string))
	}

	return Response{
		Answer:      result["answer"].(string),
		Coordinates: coordinates,
		Cells:       cells,
		Aggregator:  aggregator,
	}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	token := os.Getenv("HUGGINGFACE_TOKEN")

	csvData, err := os.ReadFile("data-series.csv")
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	table, err := CsvToSlice(string(csvData))
	if err != nil {
		fmt.Println("Error parsing CSV:", err)
		return
	}

	tapasConnector := &AIModelConnector{
		Client: &http.Client{},
	}

	fmt.Println("Welcome to Smart Home Energy Management System CLI")
	fmt.Println("-------------------------------------------------")
	fmt.Println("You can ask questions about the data in the CSV file.")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nEnter your question (or type 'q' to quit): ")
		scanner.Scan()
		input := strings.ToLower(scanner.Text())

		if input == "q" {
			fmt.Println("Exiting...")
			break
		}

		response := ProcessUserInput(input, table, tapasConnector, token)
		fmt.Println("Answer:", response)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading standard input:", err)
	}
}

func ProcessUserInput(input string, table map[string][]string, tapasConnector *AIModelConnector, token string) string {
	if strings.Contains(input, "(ask)") || strings.Contains(input, "rekomendasi") || strings.Contains(input, "recommend") {
		mistralaiConnector := &MistralAIModelConnector{
			Client: &http.Client{},
			Token:  token,
		}
		recommendation, err := mistralaiConnector.GenerateRecommendation(input)
		if err != nil {
			return fmt.Sprintf("Error generating recommendation: %v", err)
		}
		return recommendation
	}

	payload := Inputs{
		Table: table,
		Query: input,
	}
	response, err := tapasConnector.ConnectAIModel(payload, token)
	if err != nil {
		return fmt.Sprintf("Error connecting to AI model: %v", err)
	}

	return fmt.Sprintf("%s\nCoordinates: %v\nCells: %v\nAggregator: %s",
		response.Answer, response.Coordinates, response.Cells, response.Aggregator)
}
