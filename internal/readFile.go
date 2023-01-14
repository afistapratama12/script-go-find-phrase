package internal

import "os"

func ReadFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func checkPath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if _, err = os.Create(path); err != nil {
			return err
		}
	}

	return nil
}

func WriteFile(data []string, path string) error {
	if err := checkPath(path); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		return err
	}

	defer f.Close()

	for _, line := range data {
		_, err = f.WriteString(line + "\n")

		if err != nil {
			return err
		}
	}

	return nil

}

func WriteOneLine(data string, path string) error {
	if err := checkPath(path); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(data + "\n")

	if err != nil {
		return err
	}

	return nil
}
