package script

// Script represents a vulnerablity-checking script provided by the server
type Script struct {
	Name        string
	Description string
	Points      int
	Shell       string
	Script      string
}
