package infrastructure

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/jimlawless/whereami"
// )

// type Tech struct {
// 	StateFunc []func() (string, error)
// 	App       Application
// }

// type Application struct {
// 	Name             string
// 	Version          string
// 	MasterClusterUrl string
// 	Clusters         []string
// }

// func (tech *Tech) State(w http.ResponseWriter, r *http.Request) {
// 	response := TechReq{
// 		ServiceState: 0,
// 	}

// 	for _, techFunc := range tech.StateFunc {
// 		bdType, errState := techFunc()
// 		if errState != nil {
// 			response = TechReq{
// 				ServiceState: 1,
// 			}
// 			response.DbStates = append(response.DbStates, DbState{
// 				DbName:            bdType,
// 				DbState:           1,
// 				DbExceptionString: errState.Error(),
// 			})

// 			bin, err := json.Marshal(&response)
// 			if err != nil {
// 				tech.Logger.WriteError(err.Error(), whereami.WhereAmI())
// 				w.WriteHeader(500)

// 				return
// 			}

// 			w.WriteHeader(500)

// 			_, err = w.Write(bin)
// 			if err != nil {
// 				tech.Logger.WriteError(err.Error(), whereami.WhereAmI())
// 			}

// 			return
// 		}

// 		response.DbStates = append(response.DbStates, DbState{
// 			DbName:            bdType,
// 			DbState:           0,
// 			DbExceptionString: "null",
// 		})
// 	}

// 	bin, err := json.Marshal(&response)
// 	if err != nil {
// 		tech.Logger.WriteError(err.Error(), whereami.WhereAmI())
// 		w.WriteHeader(500)

// 		return
// 	}

// 	_, err = w.Write(bin)
// 	if err != nil {
// 		tech.Logger.WriteError(err.Error(), whereami.WhereAmI())
// 		w.WriteHeader(500)

// 		return
// 	}
// }

// func (tech *Tech) Info(w http.ResponseWriter, r *http.Request) {
// 	info := DbInfo{
// 		Name:    tech.App.Name,
// 		Version: tech.App.Version,
// 	}

// 	bin, err := json.Marshal(info)
// 	if err != nil {
// 		tech.Logger.WriteError(err.Error(), whereami.WhereAmI())
// 		w.WriteHeader(500)

// 		return
// 	}

// 	_, err = w.Write(bin)
// 	if err != nil {
// 		tech.Logger.WriteError(err.Error(), whereami.WhereAmI())
// 		w.WriteHeader(500)
// 	}
// }

type TechReq struct {
	ServiceState int       `json:"serviceState"`
	DbStates     []DbState `json:"dbStates"`
}

type DbState struct {
	DbName            string `json:"dbName"`
	DbState           int    `json:"dbState"`
	DbExceptionString string `json:"dbExceptionString"`
}

type DbInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
