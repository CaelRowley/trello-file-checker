# Trello File Checker
This script is a tool created in Go to query multiple Trello boards and generate a CSV file containing information about all attachments on all cards across all lists on those boards.

## Requirements

To use this script, you need to set up the following environment variables:

- `TRELLO_API_KEY`: Your Trello API key. You can obtain one by creating a new Power-Up for your workspace on the Trello [power-ups admin page](https://trello.com/power-ups/admin).
- `TRELLO_TOKEN`: Your Trello API token. You can generate a token by following the instructions provided after obtaining your API key. For example to create a read only token that expires after one month you can use the following url: 
```
https://trello.com/1/authorize?scope=read&response_type=token&key=${api_key}
```
- `TRELLO_BOARD_IDS`: A comma-separated list of Trello board IDs that you want to query. You can find the board ID in the URL of the Trello board.

## Usage

1. Clone this repository to your local machine.

2. Set up the required environment variables in your environment.

   Example:
   ```
   export TRELLO_API_KEY=your_trello_api_key
   export TRELLO_TOKEN=your_trello_token
   export TRELLO_BOARD_IDS=board_id_1,board_id_2,board_id_3
   ```

3. Build the Go executable by running:
   ```
   go build
   ```

4. Run the executable:
   ```
   ./trello-file-checker
   ```

5. Once the script finishes running, it will generate a CSV file named `attachments.csv` in the same directory as the executable. This CSV file will contain the following columns:

   - Board: The name of the Trello board.
   - List: The name of the list containing the card.
   - Card Name: The name of the card.
   - Card ID: The ID of the card.
   - File: The name of the attached file.
   - Date: The date when the attachment was added to the card.

## Notes

- This script utilizes the Trello API to fetch data. Ensure that you have the necessary permissions to access the boards you specify.
- Make sure to keep your Trello API key and token secure and do not expose them publicly.
- This script is provided as-is without any guarantees. Use it at your own risk.
