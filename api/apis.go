package api

import (
	"encoding/json"
	"errors"
	"github.com/matryer/respond"
	"net/http"
	"test_task/service"
	"test_task/store"
)

func AddDataHandler(w http.ResponseWriter, r *http.Request) {

	req := store.Data{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = errors.New("unable to get payload")
	}
	defer r.Body.Close()

	// using only basic validator due to time constraint
	if err := req.Validate(); err != nil {
		err := map[string]interface{}{"message": err.Error()}
		w.Header().Set("Content-type", "applciation/json")
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	err := service.AddData(req)
	if err != nil {
		err := map[string]interface{}{"message": err.Error()}
		w.Header().Set("Content-type", "applciation/json")
		respond.With(w, r, http.StatusConflict, err)
	} else {
		w.Header().Set("Content-type", "applciation/json")
		respond.With(w, r, http.StatusCreated, "ok")
		return
	}
}
