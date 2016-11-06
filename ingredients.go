package brewery

// Ingredientslist contains brews that could run in paralell and the order they should be run
type Ingredientslist struct {
	Brews []string         `json:"brews,omitempty"`
	Next  *Ingredientslist `json:"next,omitempty"`
}
