package f1

import "os"

func dirExists(fp string) bool {
	stat, err := os.Stat(fp)
	return err == nil && stat.IsDir()
}

func fileExists(fp string) bool {
	stat, err := os.Stat(fp)
	return err == nil && stat.IsDir()
}
