# Smart Home Energy Management System CLI

Welcome to the **Smart Home Energy Management System CLI**! This tool allows you to interact with data in a CSV file, asking questions and getting recommendations using AI models such as Google's TAPAS model and Mistral AI's Mixtral model via the Hugging Face API.

## Features

- **CSV Parsing**: Reads and processes data from a CSV file to provide insights.
- **AI-Powered Answers**: Uses Google's TAPAS model to answer questions about tabular data.
- **AI Recommendations**: Provides recommendations using Mistral AI's Mixtral model.

## Requirements

- Go 1.18 or higher
- Hugging Face API Token

## Installation

1. **Clone the repository**:

    ```bash
    git clone https://github.com/yourusername/smart-home-energy-management-cli.git
    cd smart-home-energy-management-cli
    ```

2. **Install dependencies**:

    Make sure Go modules are enabled, then run:

    ```bash
    go mod tidy
    ```

3. **Create an `.env` file**:

    Create a `.env` file in the root directory based on the `.env.example` file and provide your API tokens.

    ```plaintext
    # .env file
    HUGGINGFACE_TOKEN=your_huggingface_token_here
    ```

4. **Run the application**:

    ```bash
    go run main.go
    ```

## Usage

1. Upon running the application, you will be greeted with:

    ```plaintext
    Welcome to Smart Home Energy Management System CLI
    -------------------------------------------------
    You can ask questions about the data in the CSV file.
    ```

2. **Ask Questions**:

   You can type in your questions about the data available in the CSV file. The AI will process your question and provide answers based on the content of the CSV data.

   Example input:

    ```plaintext
    Enter your question (or type 'q' to quit): What is the average energy consumption?
    ```

3. **Get Recommendations**:

   To get recommendations, you can use phrases like `recommend`, `rekomendasi`, or `(ask)`. The application will use the Mistral AI model to generate relevant recommendations.

    ```plaintext
    Enter your question (or type 'q' to quit): recommend energy-saving tips
    ```

4. **Quit the Program**:

   Type `q` and press Enter to exit the application.

    ```plaintext
    Enter your question (or type 'q' to quit): q
    Exiting...
    ```

## Configuration

- **CSV File**: The data for processing is read from `data-series.csv`. Ensure that this file exists and is correctly formatted for the AI to understand it.
- **Environment Variables**: The application reads the Hugging Face API token from the `.env` file. Ensure that the `.env` file is properly configured with your tokens.

## Error Handling

- If the AI models return errors or the CSV file is not properly read, the CLI will display relevant error messages to help diagnose issues.
  
## Example Output

When querying the TAPAS model:

```plaintext
Enter your question (or type 'q' to quit): What is the total energy consumption?
Answer: The total energy consumption is 12345 kWh.
Coordinates: [[1, 2], [3, 4]]
Cells: ["12345", "kWh"]
Aggregator: SUM
```

When asking for recommendations with the Mistral model:

```plaintext
Enter your question (or type 'q' to quit): recommend ways to reduce energy consumption
Answer: To reduce energy consumption, consider switching to energy-efficient appliances, using smart thermostats, and turning off devices when not in use.
```

## Troubleshooting

- **API Issues**: Ensure that your Hugging Face tokens are correctly set in the `.env` file and that your API quota is not exceeded.
- **CSV Errors**: Verify that `data-series.csv` is in the correct format and accessible by the program.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request to improve the project.

## License

This project is licensed under the MIT License.
