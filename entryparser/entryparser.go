package entryparser

import (
	"fmt"
	"log"
	"net/http"

	"github.com/edkvm/sherbet/engine"
)

func output(rawCmd []string, w http.ResponseWriter) {

	cmd, err := engine.BuildCmdChain(rawCmd)

	if err != nil {
		log.Println("err_msg=", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	result, err := engine.Execute(cmd)

	if err != nil {
		log.Println("err_msg=", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Status", fmt.Sprintf("%d %s", http.StatusOK, http.StatusText(http.StatusOK)))
	if _, err := w.Write(result); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

}
