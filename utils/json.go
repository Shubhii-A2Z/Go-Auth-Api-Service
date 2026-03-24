package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJsonResponse(w http.ResponseWriter,statusCode int,data any) error{
	// Set the content type to application/json
	w.Header().Set("Content-Type","application/json"); 

	// Set the HTTP status code
	w.WriteHeader(statusCode); 

	// Encode the data as json and write it to response
	return json.NewEncoder(w).Encode(data);
}

func WriteJsonSuccessResponse(w http.ResponseWriter,statusCode int,message string,data any) error{
	response:=map[string]any{
		"message": message,
		"data": data,
		"success": true,
		"error": nil,
	}
	return WriteJsonResponse(w,statusCode,response)
}

func WriteJsonErrorResponse(w http.ResponseWriter,statusCode int,message string) error{
	response:=map[string]any{
		"message": message,
		"data": nil,
		"success": false,
	}
	return WriteJsonResponse(w,statusCode,response)
}

func ReadJsonBody(r *http.Request,result any) error{
	decoder:=json.NewDecoder(r.Body)
	return decoder.Decode(result)
}