package db

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func DBReq(name string, params map[string]any) ([](map[string]any), error) {
	dbMasterUrl, err := url.Parse(os.Getenv("DBMASTER_ORIGIN"))
	if err != nil {
		return nil, err
	}
	dbMasterUrl.Path = "/func"

	bodyJson, err := json.Marshal(map[string]any{
		"name":   name,
		"params": params,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", dbMasterUrl.String(), bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", os.Getenv("DBMASTER_KEY"))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	objectJsons := strings.Split(string(body), "\n")

	var results []map[string]any
	for _, json_ := range objectJsons {
		var object map[string]any
		err := json.Unmarshal([]byte(json_), &object)
		if err == nil {
			results = append(results, object)
		}
	}

	return results, nil
}

func GetUserData(provider string, providerId string) ([](map[string]any), error) {
	return DBReq("user.user-data", map[string]any{"provider": provider, "providerId": providerId})
}

func GetUserDataByUUID(UUID string) ([](map[string]any), error) {
	return DBReq("user.user-data-by-uuid", map[string]any{"UUID": UUID})
}

func GetUserDataByProvider(provider string, providerId string) ([](map[string]any), error) {
	return DBReq("user.user-data-by-provider-id", map[string]any{"provider": provider, "providerId": providerId})
}

func GetFileDataByFileName(fileName string) ([](map[string]any), error) {
	return DBReq("file.get-by-file-name", map[string]any{"fileName": fileName})
}

func NewFileLog(UUID string, originalFileName string, fileName string) ([](map[string]any), error) {
	return DBReq("file.log", map[string]any{"UUID": UUID, "originalFileName": originalFileName, "fileName": fileName})
}
