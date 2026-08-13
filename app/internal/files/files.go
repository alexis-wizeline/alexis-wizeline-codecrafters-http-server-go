package files

import "os"

func ReadFile(root, fileName string) ([]byte, error) {
	filePath := root + fileName
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}
