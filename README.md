# GST - Genomics Sample Tracker

GST is a tool for tracking and visualizing genomics sample processing data.
It provides functionality to export sample data to TSV files and serves an
interactive web dashboard for data exploration.

## Setup

### Environment Variables

GST requires a `.env` file in the root directory with the following variables:

```
GST_SQL_USER=gst_reader
GST_SQL_PASS=secret_password
GST_SQL_HOST=mysql.example.com
GST_SQL_PORT=3306
GST_SQL_DB=sample_tracking
```

### Building

```bash
go install
```

## Usage
GST provides two main subcommands: export and server.

### Export Subcommand
Exports sample data to a TSV file.

```
gst export --output samples.tsv
```

### Server Subcommand

```
# Start the server with default settings (port 8080)
gst server

# Specify a custom port
gst server --port 3000

# Use mock data instead of querying the database
gst server --mock ./testdata/sample_data.tsv
```

After starting the server, open your web browser and navigate to:

http://localhost:8080

The web dashboard provides:

1. A tabular view of all sample data with key information
2. A stacked horizontal bar chart showing Library Time and Sequencing Time for each sample

The chart is interactive - hover over bars to see detailed information about each sample.

#### Notes
- The database query may take several minutes to complete
- For testing purposes, use the --mock flag with a sample TSV file