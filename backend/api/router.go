package api

import (
	"context"
	"encoding/json"
	"errors"
	"events/backend/api/types"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func MagicRouter(service interface{}) http.Handler {
	mux := http.NewServeMux()
	err := register(mux, service)
	if err != nil {
		panic(err)
	}
	return mux
}

func register(mux *http.ServeMux, service interface{}) error {
	serviceType := reflect.TypeOf(service)
	serviceValue := reflect.ValueOf(service)

	if serviceType.Kind() != reflect.Ptr || serviceType.Elem().Kind() != reflect.Struct {
		return errors.New("service must be a pointer to a struct")
	}

	for i := 0; i < serviceType.NumMethod(); i++ {
		method := serviceType.Method(i)
		methodName := strings.ToLower(method.Name)

		// Ensure method signature matches: func(ctx context.Context, req T) (T, error)
		if method.Type.NumIn() != 3 || method.Type.NumOut() != 2 {
			log.Printf("WARNING: method %s has invalid number of arguments", methodName)
			continue
		}

		// Validate input arguments
		if method.Type.In(1) != reflect.TypeOf((*context.Context)(nil)).Elem() {
			log.Printf("WARNING: method %s does not take context.Context as the first argument", methodName)
			continue
		}
		if method.Type.In(2).Kind() != reflect.Struct {
			log.Printf("WARNING: method %s does not take struct as the second argument", methodName)
			continue
		}

		// Validate output values
		if method.Type.Out(0).Kind() != reflect.Pointer || method.Type.Out(0).Elem().Kind() != reflect.Struct {
			log.Printf("WARNING: method %s's first return type is not a pointer to struct", methodName)
			continue
		}
		if method.Type.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			log.Printf("WARNING: method %s does not have a second return type as an error", methodName)
			continue
		}

		// Register the HTTP handler
		log.Printf("Registering method %s", methodName)
		mux.HandleFunc("/"+methodName, func(w http.ResponseWriter, req *http.Request) {
			log.Println("Handling request for", req.URL.Path)
			if req.Method != http.MethodPost {
				errorJson(w, "only POST method is supported", http.StatusMethodNotAllowed)
				return
			}

			// Decode request body into the input struct
			reqType := method.Type.In(2)
			reqStruct := reflect.New(reqType).Interface()
			if err := json.NewDecoder(req.Body).Decode(reqStruct); err != nil {
				errorJson(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
				return
			}

			// Call the service method
			results := method.Func.Call([]reflect.Value{
				serviceValue,
				reflect.ValueOf(req.Context()),
				reflect.ValueOf(reqStruct).Elem(),
			})

			// Handle response
			resValue := results[0].Interface()
			errValue := results[1].Interface()

			if errValue != nil {
				errorJson(w, errValue.(error).Error(), http.StatusInternalServerError)
				log.Printf("ERROR: %v", errValue)
				return
			}

			// Write the response
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(resValue)
			if err != nil {
				log.Println("ERROR: encoding response:", err)
			}
		})
	}

	return nil
}

func errorJson(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	e := json.NewEncoder(w).Encode(types.ErrorResponse{ErrorMessage: err})
	if e != nil {
		log.Println("ERROR: encoding error response:", e)
	}
}
