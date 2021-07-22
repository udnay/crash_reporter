package main

import "os"

func GetCorePattern() (string, error) {

	f, err := os.Open("/proc/sys/kernel/core_pattern")
	if err != nil {
		println("Unable to open core_pattern file")
	}
	defer f.Close()

	b := make([]byte, 2048)

	_, err = f.Read(b)
	if err != nil {
		println("Unable to read file")
	}
	return string(b), nil
}

func SetCorePattern() error {
	f, err := os.OpenFile("/proc/sys/kernel/core_pattern", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		print("unable to open core_pattern")
		return err
	}
	defer f.Close()

	core_pattern := "|/path/to/my/executable"
	_, err = f.Write([]byte(core_pattern))
	if err != nil {
		print("Unable to write core_pattern")
		return err
	}

	return nil
}
