package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

const baseURL = "https://mage.uibot.com.cn/v1"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TaskID  string      `json:"task_id,omitempty"`
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func getHeaders(publicKey string, secretKey string) http.Header {
	if publicKey == "" || secretKey == "" {
		log.Fatalln("Invalid public key or secret key")
	}
	var nounce, timestampString, signToken string
	var err error

	nounce, err = randomHex(16)
	if err != nil {
		log.Fatalln(err)
	}

	timestampString = strconv.FormatInt(time.Now().Unix(), 10)

	h := sha1.New()
	h.Write([]byte(nounce + timestampString + secretKey))
	signToken = hex.EncodeToString(h.Sum(nil))

	return http.Header{
		"Api-Auth-nonce":     {nounce},
		"Api-Auth-pubkey":    {publicKey},
		"Api-Auth-timestamp": {timestampString},
		"Api-Auth-sign":      {signToken},
		"Content-Type":       {"application/json"},
	}
}

func post(endPoint string, params interface{}, headers http.Header, result *Response) error {
	jsonValue, err := json.Marshal(params)
	if err != nil {
		return err
	}
	u, _ := url.Parse(baseURL)
	u.Path = path.Join(u.Path, endPoint)
	url := u.String()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	req.Header = headers

	c := &http.Client{Timeout: 30 * time.Second}
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// use either way to decode response
	// 1. json decoder decode stream to struct (more efficient)
	// 2. use ioutil.ReadAll if you want to see server's native response for debug

	// 1
	json.NewDecoder(res.Body).Decode(result)

	// 2
	// resBody, _ := ioutil.ReadAll(res.Body)
	// fmt.Printf("Server Response: \n %s\n\n", string(resBody))
	// err = json.Unmarshal(resBody, result)

	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func readFileBase64(filePath string) string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalln(err)
	}
	return base64.StdEncoding.EncodeToString(content)
}

// ==================================
// NLP
// ==================================
func NormalizeAddress(address string, result *Response) error {
	params := map[string]string{"text": address}
	return post("/mage/nlp/geoextract", params, getHeaders(GeoExtractPublicKey, GeoExtractSecretKey), result)
}

// ????????????
func ClassifyDocumentContent(doc string, result *Response) error {
	params := map[string]string{"doc": doc}
	return post("/document/classify", params, getHeaders(DocContentClassifyPublicKey, DocContentClassifySecretKey), result)
}

// ????????????
func GetMatchText(text string, result *Response) error {
	params := map[string]string{"text": text}
	return post("/mage/nlp/textmatch", params, getHeaders(TextMatchPublicKey, TextMatchSecretKey), result)
}

// ????????????
func ExtractInfoFromDocumentContent(docContent string, result *Response) error {
	params := map[string]string{"doc": docContent}
	return post("/document/extract", params, getHeaders(DocContentExtractPublicKey, DocContentExtractSecretKey), result)
}

// ????????????-????????????
// This api will return a task_id, you can use this task_id to query the result later on
func SubmitDocument(filePath string, result *Response) error {
	params := map[string]string{"file_base64": readFileBase64(filePath)}
	return post("/mage/nlp/docextract/create", params, getHeaders(DocExtractPublicKey, DocExtractSecretKey), result)
}

// ????????????-????????????
// Use this api to get the process result with the task_id given from ????????????-???????????? api
func QueryDocumentResult(taskID string, result *Response) error {
	params := map[string]string{"task_id": taskID}
	return post("/mage/nlp/docextract/query", params, getHeaders(DocExtractPublicKey, DocExtractSecretKey), result)
}

// ==================================
// OCR
// ==================================

// ???????????????
func OcrCaptcha(filePath string, result *Response) error {
	params := map[string]string{"img_base64": readFileBase64(filePath)}
	return post("/document/ocr/verification", params, getHeaders(VerificationPublicKey, VerificationSecretKey), result)
}

// ??????????????????
func OcrLicense(filePath string, result *Response) error {
	params := map[string]string{"img_base64": readFileBase64(filePath)}
	return post("/document/ocr/license", params, getHeaders(LicensePublicKey, LicenseSecretKey), result)
}

// ????????????
func OcrStamp(filePath string, result *Response) error {
	params := map[string]string{"img_base64": readFileBase64(filePath)}
	return post("/document/ocr/stamp", params, getHeaders(StampPublicKey, StampSecretKey), result)
}

// ?????????????????????
func OcrBill(filePath string, result *Response) error {
	params := map[string]string{"img_base64": readFileBase64(filePath)}
	return post("/document/ocr/bills", params, getHeaders(BillsPublicKey, BillsSecretKey), result)
}

// ??????????????????
func OcrTable(filePath string, result *Response) error {
	params := map[string]interface{}{
		"img_base64": []string{readFileBase64(filePath)},
	}
	return post("/document/ocr/table", params, getHeaders(TablePublicKey, TableSecretKey), result)
}

// ????????????
func OcrTemplate(filePath string, result *Response) error {
	params := map[string]string{"img_base64": readFileBase64(filePath)}
	return post("/document/ocr/template", params, getHeaders(TemplatePublicKey, TemplateSecretKey), result)
}

// ??????????????????
func OcrGeneral(filePath string, result *Response) error {
	params := map[string]interface{}{
		"img_base64": []string{readFileBase64(filePath)},
	}
	return post("/document/ocr/general", params, getHeaders(GeneralPublicKey, GeneralSecretKey), result)
}

// ==================================
// Contract
// ==================================
func SubmitContract(filePath1 string, filePath2 string, result *Response) error {
	params := map[string]string{
		"file_base":    readFileBase64(filePath1),
		"file_compare": readFileBase64(filePath2),
	}
	return post("/mage/solution/contract/compare", params, getHeaders(ContractPublicKey, ContractSecretKey), result)
}
func QueryContractResult(taskID string, result *Response) error {
	params := map[string]string{"task_id": taskID}
	return post("/mage/solution/contract/detail", params, getHeaders(ContractPublicKey, ContractSecretKey), result)
}

// get download link
func DownloadContract(taskID string, result *Response) error {
	params := map[string]string{"task_id": taskID}
	return post("/mage/solution/contract/files", params, getHeaders(ContractPublicKey, ContractSecretKey), result)
}

// ==================================
// IDP
// ==================================
func SubmitFlow(filePath string, fileName string, result *Response) error {
	params := make(map[string]map[string]string)
	params["file"] = map[string]string{
		"base64": readFileBase64(filePath),
		"name":   fileName,
	}
	return post("/mage/idp/flow/task/create", params, getHeaders(FlowPublicKey, FlowSecretKey), result)
}
func QueryFlowResult(taskID string, withOcrResult bool, result *Response) error {
	params := make(map[string]interface{})
	params["task_id"] = taskID
	params["with_ocr_general"] = withOcrResult
	return post("/mage/idp/flow/task/query", params, getHeaders(ContractPublicKey, ContractSecretKey), result)
}
