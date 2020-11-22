package transport

type errorResponse struct {
	Error     string `json:"error"`
	ErrorCode uint   `json:"error_code"`
}

type BuildingsDictionaryResposnse struct {
	Data []*Building `json:"data"`
}

type EntriesDictionaryResponse struct {
	Data []*Entry `json:"data"`
}

type Building struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Entry struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
