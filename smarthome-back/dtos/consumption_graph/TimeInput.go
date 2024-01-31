package consumption_graph

type TimeInput struct {
	Type            string   `json:"type" binding:"required"`
	SelectedOptions []string `json:"selectedOptions" binding:"required"`
	Time            string   `json:"time" binding:"required"`
	QueryType       string   `json:"queryType" binding:"-"`
}
