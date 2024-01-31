package consumption_graph

type DateInput struct {
	Type            string   `json:"type" binding:"required"`
	SelectedOptions []string `json:"selectedOptions" binding:"required"`
	Start           string   `json:"start" binding:"required"`
	End             string   `json:"end" binding:"required"`
	QueryType       string   `json:"queryType" binding:"-"`
}
