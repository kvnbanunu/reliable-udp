package utils

import "os"

// Creates logfile if not exists
func PrepareLogFile(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o777)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

// opens file for reading only
func OpenLogFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY, 0o777)
}

// Writes new data to a temporary file, then replaces the real log file
func AtomicWrite(logDir, prog, logPath string, data []byte) error {
	tempFile, err := os.CreateTemp(logDir, "tempfile_"+prog)
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(data)
	if err != nil {
		return err
	}
	tempFile.Close()

	err = os.Rename(tempFile.Name(), logPath)
	if err != nil {
		return err
	}
	return nil
}
