package main

func getWd() (string, error) {
	// Placeholder implementation: In a real application, this function would return the current working directory.
	// For now, it returns a static path for demonstration purposes.
	return "/home/frank", nil
}

func readDir(path string) ([]string, error) {
	// Placeholder implementation: In a real application, this function would read the directory contents.
	// For now, it returns a static list of items for demonstration purposes.
	return []string{"file1.txt", "file2.txt", "file3.txt"}, nil
}