package common

// Method ...
type Method struct {
	Name      string     `json:"name"`
	Arguments []Argument `json:"arguments"`
}

// Argument ...
type Argument struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
