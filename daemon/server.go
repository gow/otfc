package main

import (
	"encoding/json"
	"fmt"
	"github.com/gow/otfc/config"
	"log"
	"net/http"
)

const (
	HTTP_PORT = "8088"
)

type httpServer struct {
	configPtr *config.ConfigFile
}

func (server *httpServer) start() error {
	http.HandleFunc(
		"/set",
		func(w http.ResponseWriter, r *http.Request) {
			server.httpCallbackSet(w, r)
		})
	http.HandleFunc(
		"/get",
		func(w http.ResponseWriter, r *http.Request) {
			server.httpCallbackGet(w, r)
		})
	http.HandleFunc(
		"/delete",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println(server)
			log.Println("Handling delete request")
			server.httpCallbackDelete(w, r)
		})
	err := http.ListenAndServe(":"+HTTP_PORT, nil)
	if err != nil {
		return err
	}
	return err
}

func (server *httpServer) httpCallbackGet(
	resp http.ResponseWriter,
	req *http.Request) {

	key := req.URL.Query().Get("key")
	value, err := server.configPtr.Get(key)
	if err != nil {
		sendHttpError(resp, err, http.StatusBadRequest)
		return
	}
	sendHttpJSONResponse(
		resp,
		struct {
			Status string
			Key    string
			Value  []byte
		}{"OK", key, value})
}

func (server *httpServer) httpCallbackSet(
	resp http.ResponseWriter,
	req *http.Request) {

	key := req.URL.Query().Get("key")
	value := req.URL.Query().Get("value")
	log.Println("Key: ", key, "Value: ", value)
	if value == "" {
		sendHttpError(
			resp,
			Error{ErrNo: ERR_DMN_INVALID_VALUE},
			http.StatusNotAcceptable)
		return
	}
	err := server.configPtr.Set(key, []byte(value))
	if err != nil {
		sendHttpJSONResponse(resp, err)
		return
	}
	sendHttpJSONResponse(
		resp,
		struct {
			Status string
			Key    string
			Value  string
		}{"OK", key, value})
}

func (server *httpServer) httpCallbackDelete(
	resp http.ResponseWriter,
	req *http.Request) {

	key := req.URL.Query().Get("key")
	err := server.configPtr.Delete(key)
	log.Println("httpCallbackDelete [key, err]: ", key, err)
	if err != nil {
		sendHttpError(resp, err, http.StatusBadRequest)
		return
	}
	sendHttpJSONResponse(
		resp,
		struct {
			Status string
			Key    string
		}{"OK", key})
}

func sendHttpError(w http.ResponseWriter, err interface{}, errCode int) {
	response := struct {
		Status string
		Err    interface{}
	}{"error", err}
	jsonResponse, _ := json.Marshal(response)
	http.Error(w, string(jsonResponse)+"\n", errCode)
}

func sendHttpJSONResponse(w http.ResponseWriter, data interface{}) {
	jsonResponse, _ := json.Marshal(data)
	fmt.Fprintf(w, string(jsonResponse)+"\n")
}
