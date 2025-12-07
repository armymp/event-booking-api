# How to Launch the Event Booking API
### This project uses Go, Gin, SQLite, and Viper for environment-based configuration.
### You can run the API in **development** or **production** mode.

## Prerequisites
1. Install Go (version 1.20 or later recommended)
2. Make sure your Go environment variables are set up (`GOPATH` and `PATH`)
3. Install dependencies: run `go mod tidy` in the project root
4. if using VS Code and wanting to debug, install the Delve debugger:
    - Run: `go install -v github.com/go-delve/delve/cmd/dlv@latest`
    - Make sure your Go bin directory is in your PATH

<br>

## Configuration
The app uses environment-based configuration using `.env` files and Viper.
- `.env` - development configuration (default)
    - Example contents:
        APP_ENV=development
        DB_NAME=dev.db
- `.env.production` - production configuration
    - Example contents:
        APP_ENV=production
        DB_NAME=events.db

The `config.LoadConfig()` function reads these files automatically based on the `APP_ENV` environment variable.

<br>

## Running in Devlopment Mode
Option 1: Using VS Code
1. Open the **Run and Debug** panel.
2. Select **Run Dev API**.
3. Click the green play button

Option 2: Using Terminal / Command Line
- On Windows Powershell:
    `$env:APP_ENV="development"`
    `go run main.go`

- on macOS/Linux:
    `APP_ENV=development go run main.go`

By default, the development server runs on `http://localhost:8080` and uses the database `dev.db`.

<br>

## Running in Production Mode
Option 1: Using VS Code
1. Open the **Run and Debug** panel.
2. Select **Run Prod API**.
3. Click the green play button

Option 2: Using Terminal / Command Line
- On Windows Powershell:
    `$env:APP_ENV="production"`
    `go run main.go`

- on macOS/Linux:
    `APP_ENV=production go run main.go`

By default, the development server runs on `http://localhost:80` and uses the database `events.db`.

<br>

### Troubleshooting
- Forking this repo automatically gives you both developoment and production configurations
- You do not need to modify the code to switch environments
- All environment-specific settings are in `.env` files or `launch.json` for VS Code
- Make sure to run `go mod tidy` before starting to ensure all dependencies are installed
