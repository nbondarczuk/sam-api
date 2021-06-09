package common

import (
	"bytes"
    "compress/gzip"	
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/gorilla/mux"
)

type Any interface{}

type QueryFilter map[string]Any

// dump what was receivedfrom the request
func RequestInfo(r *http.Request) string {
	dump, _ := httputil.DumpRequest(r, true)
	info := fmt.Sprintf("Request %s:%s %s:%s %T %q",
		"method", r.Method,
		"url", r.URL.Path,
		dump, dump)

	return info
}

// Show request heders
func LogHeadersInfo(r *http.Request) {
	for name, headers := range r.Header {
		for _, h := range headers {
			log.Printf("%v: %v\n", name, h)
		}
	}

	return
}

// Get query parameters
func UrlQueryParam(r *http.Request, key string, optional bool) (value string, err error) {
	keys, ok := r.URL.Query()[key]
	value = keys[0]
	if !ok || value == "" {
		if optional {
			return
		} else {
			return value, fmt.Errorf("Manatatory url parameter missing: " + key)
		}
	}

	return
}

// Get & validate existence of string parameter
func PathVariableStr(r *http.Request, label string, mandatory bool) (value string, err error) {
	vars := mux.Vars(r)
	value, ok := vars[label]
	if !ok {
		if mandatory {
			err = fmt.Errorf("Manadatory variable does not exist: " + label)
		}
	}

	return value, err
}

// Get & validate existence of int64 parameter
func PathVariableInt64(r *http.Request, label string, mandatory bool) (value int64, err error) {
	vars := mux.Vars(r)
	v, ok := vars[label]
	if !ok {
		if mandatory {
			err = fmt.Errorf("Manadatory variable does not exist: " + label)
		}
	} else {
		value, err = strconv.ParseInt(v, 10, 64)
	}

	return value, err
}

func ReadPayload(r *http.Request) (body []byte, err error) {
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Can't read request body")
	}
	
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return
}

func gunzip(body []byte) (decobody []byte, err error) {
	buf := bytes.NewBuffer(body)
	
    reader, err := gzip.NewReader(buf)
    if err != nil {
        return
    }
    defer reader.Close()

    decobody, err = ioutil.ReadAll(reader)
    if err != nil {
        return
    }	

	return
}

func GetPayload(r *http.Request, encoding string) (body []byte, err error) {
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Can't read request body")
	}

	if encoding == "gzip" {
		if body, err = gunzip(body); err != nil {
			return nil, fmt.Errorf("Failed to gunzip payload: %s", err.Error())
		}
	}

	return
}

func GetAttributesWithValues(r *http.Request) (attributes *map[string]interface{}, body []byte, err error) {
	// get body of the request
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("Can't read request body")
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	log.Printf("Parsing dynamic json: %s", string(body))

	// get attributes and their values
	var result, data map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	if result["data"] != nil {
		data = result["data"].(map[string]interface{})
	} else {
		data = make(map[string]interface{})
	}
	
	log.Printf("Dynamic json loaded: %#v", data)
	
	return &data, body, nil
}

// Get & validate existence of all attributes in data json record
func GetEvaluations(r *http.Request) (data *map[string]interface{}, body []byte, err error) {
	if data, body, err = GetAttributesWithValues(r); err != nil {
		return nil, nil, err
	} 

	log.Printf("Loaded evaluations: %#v", *data)
	
	return data, body, nil
}

func EmptyValue(v interface{}) (empty bool) {
	switch v.(type) {
	case string:
		if v.(string) == "" {
			empty = true
		}
	default:
		empty = false
	}

	return
}

func MemberOf(v interface{}, values ...string) (ok bool) {
	ok = false
	switch v.(type) {
	case string:
		for _, vv := range values {
			if vv == v {
				ok = true
				break
			}
		}
	default:
		ok = false
	}

	return
}
