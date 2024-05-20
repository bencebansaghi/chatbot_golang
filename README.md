# Instaloader golang 

This Go application fetches Instagram posts from specified profiles within the last day and outputs them in JSON format. It is designed for integrating with the student-events Telegram bot, addressing issues encountered with Instagram flagging the previously used Python library.


## Features

- Fetches recent posts from specified Instagram profiles.
- Outputs post details in JSON format, including caption, shortcode, and username.
- Utilizes environment variables for configuration.
- Supports loading Instagram credentials from a `.env` file or a `.goinsta` configuration file.

## Configuration

Before running the application, configure the following environment variables in a `.env` file:

- `INSTA_USERNAME`: Your Instagram username.
- `INSTA_PASSWORD`: Your Instagram password.
- `INSTA_PROFILES`: Comma-separated list of Instagram profiles to fetch posts from.

Example `.env` file:

```plaintext
INSTA_USERNAME=your_username
INSTA_PASSWORD=your_password
INSTA_PROFILES=profile1,profile2,profile3
```

## Dependencies

- [goinsta/v3](https://github.com/Davincible/goinsta): Instagram private API.
- [joho/godotenv](https://github.com/joho/godotenv): For loading environment variables from a `.env` file.

## Installation

1. Ensure you have Go installed on your system.
2. Clone the repository and navigate into the project directory.
3. Install the required Go modules:

```shell
go mod tidy
```

## Running the Application

To run the application, execute the following command in the terminal:

```shell
go run main.go
```

The application will fetch posts from the specified Instagram profiles and output them in JSON format to the standard output.

## Output Format

Each post is outputted in the following JSON format:

```json
{
  "Caption": "Post caption here",
  "ShortCode": "Post shortcode",
  "Username": "Profile username"
}
```

### usage.py

This Python script demonstrates the example usage of how to run the Go application as a subprocess and handle its output.
