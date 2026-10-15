Create a Command Line Interface (CLI) application that acts as a file unpacker/splitter from a single text stream. 

The application must read a text file (or standard input) containing multiple file payloads separated by specific block markers, parse them, and write them to the local file system.

### 1. Marker Specification
Blocks are defined by a starting marker line. The marker format is:
##-## [FILE_PATH] ##-##

- Example: ##-## out/config.go ##-##
- Everything following this marker line, up until the next marker line or the end of the file, constitutes the raw content of that specific file.
- Leading and trailing blank lines immediately adjacent to the markers should be trimmed from the content.

### 2. Core Requirements & Logic
- **Path Parsing:** Extract the target file path from each marker. Paths can be relative (e.g., `out/config.go`) or absolute.
- **Directory Creation:** For each file, check if its parent directory structure exists. If it does not exist, create it automatically (equivalent to `mkdir -p`).
- **File Writing:** Write the extracted content to the parsed path. If the file already exists, overwrite it.
- **Error Handling:** Gracefully handle and report issues such as invalid marker syntax, permission denied errors when creating directories, or failures writing files.

### 3. Architecture & Future Extensibility
- **Separation of Concerns:** Separate the core parsing/writing logic from the CLI interface. 
- **Interface Design:** Design the core processing logic so it accepts an input stream (like an io.Reader) and a target base directory. This ensures the engine can easily be wrapped in a Web API (e.g., accepting a file upload via POST request) or a GUI in the future without rewriting the parser.
- **Current CLI UX:** For now, the CLI should accept the input file path as a command-line argument. Provide clear console output showing which files were successfully created.