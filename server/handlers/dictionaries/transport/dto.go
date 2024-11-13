package transport

import "github.com/ivch/dynasty/server/handlers/requests"

type errorResponse struct {
	Error     string `json:"error"`
	ErrorCode uint   `json:"error_code"`
	Ru        string `json:"ru"`
	Ua        string `json:"ua"`
}

type RequestTypesDictionaryResponse struct {
	Data map[requests.RequestType]map[string]string `json:"data"`
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
