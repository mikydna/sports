package livetiming

import "os"

func fileExists(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil
}
