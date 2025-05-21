
# 📘 Concurrent Log Processor

## This Go project implements a concurrent log processing tool designed to efficiently read a large log file, count the frequency of specified keywords (e.g., INFO, ERROR, DEBUG), and output them in descending order of occurrence.

## ✨ Features

✅ Reads large log files efficiently using streaming.

⚙️ Uses Goroutines and Worker Pools for concurrent line processing.

🔒 Thread-safe keyword aggregation.

📊 Outputs sorted frequency of keywords.

🧪 Includes unit test.

🧪 Includes integration test for full flow validation.

🏗️ Project Structure

<pre> ``` . 
   ├── main.go             # Main logic to process logs
   ├── processor
   │   └── log.txt         # Sample log file (99 lines)
   │   └── file_test.go    # unit test   
   │   └── file.go         # reading file 
   │   └── process_file.go # process file for concurrent logging     
   │   └── process_file_test.go #  unit test   
   │   └── log_test.go     # Integration test   
   ├── go.mod
   └── README.md
</pre>

## 📥 Requirements

Go 1.24.1
🚀 How to Run

### Clone the repository:
git clone github.com/funmi4194/log-processor
cd test_folder
Run the application:
go run main.go


### 🧪 Running Tests
Run the integration test with:
go test -v
Ensured your log.txt file contains testable keyword occurrences.
📄 Example Log Entry

2023-10-28 12:00:01 - INFO - User logged in

2023-10-28 12:00:03 - ERROR - Database connection failed

2023-10-28 12:00:04 - DEBUG - Cache hit for request
✅ Expected Output

INFO: 1

ERROR: 1

DEBUG: 1


### 🧠 Notes

Keywords are case-insensitive.
The number of workers and batch size can be configured for optimal performance.
Buffered channels are used to prevent blocking and improve throughput.


👨‍💻 Author
[Olayiwola Oluwafunmilayo]
