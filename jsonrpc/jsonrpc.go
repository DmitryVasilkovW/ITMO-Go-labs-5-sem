//go:build !solution

package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func MakeHandler(service interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		methodName := extractMethodName(r)
		method, err := findMethod(service, methodName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		reqValue, err := parseRequest(r, method)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rsp, err := callMethod(ctx, method, reqValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := writeResponse(w, rsp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func extractMethodName(r *http.Request) string {
	return r.URL.Path[1:]
}

func findMethod(service interface{}, methodName string) (reflect.Value, error) {
	method := reflect.ValueOf(service).MethodByName(methodName)
	if !method.IsValid() {
		return reflect.Value{}, errors.New("method not found")
	}
	return method, nil
}

func parseRequest(r *http.Request, method reflect.Value) (reflect.Value, error) {
	reqType := method.Type().In(1)
	reqValue := reflect.New(reqType.Elem())

	if err := json.NewDecoder(r.Body).Decode(reqValue.Interface()); err != nil {
		return reflect.Value{}, err
	}
	return reqValue, nil
}

func callMethod(ctx context.Context, method reflect.Value, reqValue reflect.Value) (interface{}, error) {
	result := method.Call([]reflect.Value{reflect.ValueOf(ctx), reqValue})

	if len(result) != 2 {
		return nil, errors.New("invalid method signature")
	}

	if !result[1].IsNil() {
		err := result[1].Interface().(error)
		return nil, err
	}

	return result[0].Interface(), nil
}

func writeResponse(w http.ResponseWriter, rsp interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(rsp)
}

func Call(ctx context.Context, endpoint string, method string, req, rsp interface{}) error {
	url := buildURL(endpoint, method)

	reqBody, err := marshalRequest(req)
	if err != nil {
		return err
	}

	httpReq, err := createRequest(ctx, url, reqBody)
	if err != nil {
		return err
	}

	resp, err := sendRequest(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := checkResponseStatus(resp); err != nil {
		return err
	}

	return decodeResponse(resp, rsp)
}

func buildURL(endpoint, method string) string {
	return fmt.Sprintf("%s/%s", endpoint, method)
}

func marshalRequest(req interface{}) ([]byte, error) {
	return json.Marshal(req)
}

func createRequest(ctx context.Context, url string, body []byte) (*http.Request, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	return httpReq, nil
}

func sendRequest(req *http.Request) (*http.Response, error) {
	httpClient := http.DefaultClient
	return httpClient.Do(req)
}

func checkResponseStatus(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(bodyBytes))
	}
	return nil
}

func decodeResponse(resp *http.Response, rsp interface{}) error {
	return json.NewDecoder(resp.Body).Decode(rsp)
}
