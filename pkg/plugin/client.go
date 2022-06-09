package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type FacetRequest struct {
	QueryType string `json:"queryType"`
	MaxCount  int    `json:"maxCount"`

	Field string `json:"field"`
}

type DataSetClient struct {
	dataSetUrl string
	apiKey     string
	netClient  *http.Client
}

func NewDataSetClient(dataSetUrl string, apiKey string) *DataSetClient {
	// Consider using the backend.httpclient package provided by the Grafana SDK.
	// This would allow a per-instance configurable timeout, rather than the hardcoded value here.
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	return &DataSetClient{
		dataSetUrl: dataSetUrl,
		apiKey:     apiKey,
		netClient:  netClient,
	}
}

func (d *DataSetClient) doPingRequest(req interface{}) (*LRQResult, error) {
	// Long-Running Query (LRQ) api usage:
	// - An initial POST request is made containing the standard/power query
	// - Its response may or may not contain the results
	//   - This is indicated by StepsCompleted == StepsTotal in the response
	// - If not complete, follow up with GET ping requests with the response Id
	//   - If the token is present in the initial POST request response, include it in subsequent pings
	// - When complete send a DELETE request to clean up resources
	//   - If the token is present in the initial POST request response, include it in this request as well

	const TOKEN_HEADER = "X-Dataset-Query-Forward-Tag"

	body, err := json.Marshal(req)
	if err != nil {
		log.DefaultLogger.Error("error marshalling request to DataSet", "err", err)
		return nil, err
	}

	request, err := http.NewRequest("POST", d.dataSetUrl+"/v2/api/queries", bytes.NewBuffer(body))
	if err != nil {
		log.DefaultLogger.Error("error constructing request to DataSet", "err", err)
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+d.apiKey)
	request.Header.Set("Content-Type", "application/json")

	var respBody LRQResult
	var token string
	for i := 0; ; i++ {
		resp, err := d.netClient.Do(request)
		if err != nil {
			if e, ok := err.(*url.Error); ok && e.Timeout() {
				log.DefaultLogger.Error("request to DataSet timed out")
				return nil, e
			} else {
				return nil, err
			}
		}

		respBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.DefaultLogger.Error("error reading response from DataSet", "err", err)
			return nil, err
		}

		if err = json.Unmarshal(respBytes, &respBody); err != nil {
			log.DefaultLogger.Error("error unmarshaling response from DataSet", "err", err)
			return nil, err
		}

		if respBody.StepsCompleted >= respBody.StepsTotal {
			break
		}

		// Only check for the token from the initial POST launch request
		if i == 0 {
			token = resp.Header.Get(TOKEN_HEADER)
		}

		time.Sleep(100 * time.Millisecond)

		u := fmt.Sprintf("%s/v2/api/queries/%s?lastStepSeen=%d", d.dataSetUrl, respBody.Id, respBody.StepsCompleted)
		request, err = http.NewRequest("GET", u, nil)
		if err != nil {
			log.DefaultLogger.Error("error constructing request to DataSet", "err", err)
			return nil, err
		}
		request.Header.Set("Authorization", "Bearer "+d.apiKey)
		request.Header.Set("Content-Type", "application/json")
		if token != "" {
			request.Header.Set(TOKEN_HEADER, token)
		}
	}

	u := fmt.Sprintf("%s/v2/api/queries/%s", d.dataSetUrl, respBody.Id)
	request, err = http.NewRequest("DELETE", u, nil)
	if err != nil {
		log.DefaultLogger.Warn("error constructing request to DataSet", "err", err)
	} else {
		request.Header.Set("Authorization", "Bearer "+d.apiKey)
		request.Header.Set("Content-Type", "application/json")
		if token != "" {
			request.Header.Set(TOKEN_HEADER, token)
		}
		if resp, err := d.netClient.Do(request); err != nil {
			if e, ok := err.(*url.Error); ok && e.Timeout() {
				log.DefaultLogger.Warn("request to DataSet timed out")
			} else {
				log.DefaultLogger.Warn("error sending request to DataSet", "err", err)
			}
		} else {
			// Read/close the body so the client's transport can re-use a persistent tcp connection
			io.ReadAll(resp.Body)
			resp.Body.Close()
		}
	}

	return &respBody, nil
}

func (d *DataSetClient) DoLRQRequest(req LRQRequest) (*LRQResult, error) {
	return d.doPingRequest(req)
}

func (d *DataSetClient) DoFacetValuesRequest(req FacetQuery) (*LRQResult, error) {
	return d.doPingRequest(req)
}

func (d *DataSetClient) DoTopFacetRequest(req TopFacetRequest) (*LRQResult, error) {
	return d.doPingRequest(req)
}

func (d *DataSetClient) DoFacetRequest(req FacetRequest) (int, error) {
	body, err := json.Marshal(req)
	if err != nil {
		log.DefaultLogger.Error("error marshalling request to DataSet", "err", err)
		return 0, err
	}

	request, err := http.NewRequest("POST", d.dataSetUrl+"/api/facetQuery", bytes.NewBuffer(body))
	if err != nil {
		log.DefaultLogger.Error("error constructing request to DataSet", "err", err)
		return 0, err
	}
	request.Header.Set("Authorization", "Bearer "+d.apiKey)
	request.Header.Set("Content-Type", "application/json")

	resp, err := d.netClient.Do(request)
	if err != nil {
		log.DefaultLogger.Error("error sending request to DataSet", "err", err)
		return 0, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.DefaultLogger.Error("error reading response from DataSet", "err", err)
		return 0, err
	}
	log.DefaultLogger.Debug("Result of request to facet", "body", string(respBytes))

	return resp.StatusCode, nil
}