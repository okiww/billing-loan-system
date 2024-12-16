package helpers

import "flag"

// IsTestEnv checks if the application is running in a test environment
func IsTestEnv() bool {
	return flag.Lookup("test.v") != nil
}
