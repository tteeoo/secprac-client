package script

// Script represents a vulnerablity-checking script provided by the server
type Script struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Points      int    `json:"points"`
	Shell       string `json:"shell"`
	Script      string `json:"script"`
}
