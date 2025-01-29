# Fit to CSV Processor

## Overview
This project provides a web-based tool for uploading `.fit` files, processing them into cleaned `.csv` files using a Python script, and downloading the processed files.

The backend is built with Go, and it executes a Python script to handle data cleaning and conversion.

## Features
- Upload `.fit` files via a web interface.
- Process the files with a Python script that converts them to `.csv`.
- Download the cleaned `.csv` files.

## Project Structure
```
fit-to-csv/
│── uploads/        # Stores uploaded .fit files
│── processed/      # Stores processed .csv files
│── static/         # Contains HTML, CSS, and JavaScript for the frontend
│── process_fit.py  # Python script to process .fit files
│── main.go         # Go backend handling uploads and downloads
│── README.md       # Project documentation
```

## Requirements
- Go 1.18+
- Python 3.x
- pandas library (`pip install pandas`)
- fitparse library (`pip install fitparse`)

## Setup
1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/fit-to-csv.git
   cd fit-to-csv
   ```
2. Install dependencies:
   - Ensure Go is installed.
   - Install Python dependencies:
     ```sh
     pip install pandas fitparse
     ```
3. Run the Go backend:
   ```sh
   go run main.go
   ```
4. Access the web interface:
   - Open `http://localhost:8080` in your browser.

## API Endpoints
- `POST /upload` - Upload a `.fit` file.
- `GET /download?file=filename.csv` - Download the processed `.csv` file.

## Python Processing Script (`process_fit.py`)
This script loads `.fit` files, extracts relevant data, processes timestamps, and outputs a cleaned `.csv` file.

### Usage:
```sh
python process_fit.py file1.fit file2.fit
```

### Selected Records:
The script extracts and processes the following fields:
- `record`
- `timestamp`
- `distance`
- `enhanced_altitude`
- `enhanced_speed`
- `gps_accuracy`
- `position_lat`
- `position_long`
- `speed`
- `heart_rate`

### Timestamp Processing:
- Converts timestamps to ISO 8601 format.
- Optionally rounds timestamps to a specified interval (default: 0 seconds).

## Future Enhancements
- Add user authentication.
- Improve UI/UX.
- Enhance error handling and logging.

## License
MIT License. See `LICENSE` for details.

