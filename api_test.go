package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeAddress(t *testing.T) {
	result := Response{}
	err := NormalizeAddress("徐汇区虹漕路70号", &result)
	assert.Nil(t, err)
	t.Log(result.Data)
}
func TestClassifyDocumentContent(t *testing.T) {
	result := Response{}
	err := ClassifyDocumentContent("西瓜", &result)
	assert.Nil(t, err)
	t.Log(result)
}
func TestGetMatchText(t *testing.T) {
	result := Response{}
	err := GetMatchText("西瓜", &result)
	assert.Nil(t, err)
	t.Log(result)
}
func TestExtractInfoFromDocumentContent(t *testing.T) {
	result := Response{}
	err := ExtractInfoFromDocumentContent("韩国通信运营商KT的有线、无线网络发生了无法连接的状况，导致全国企业、餐馆、普通家庭等上不了网，带来不便。", &result)
	assert.Nil(t, err)
	t.Log(result)
}
func TestOcrCaptcha(t *testing.T) {
	result := Response{}
	err := OcrCaptcha("./files/captcha.jpeg", &result)
	assert.Nil(t, err)
	t.Log(result)
}
func TestOcrLicense(t *testing.T) {
	result := Response{}
	err := OcrLicense("./files/id_card.jpeg", &result)
	assert.Nil(t, err)
	t.Log(result)
}
func TestOcrStamp(t *testing.T) {
	result := Response{}
	err := OcrStamp("./files/red_stamp.png", &result)
	assert.Nil(t, err)
	t.Log(result)
}
func TestOcrBill(t *testing.T) {
	result := Response{}
	err := OcrBill("./files/bill.jpg", &result)
	assert.Nil(t, err)
	t.Log(result)
}
func TestOcrTable(t *testing.T) {
	result := Response{}
	err := OcrTable("./files/table.png", &result)
	assert.Nil(t, err)
	t.Log(result)
}
func TestOcrTemplate(t *testing.T) {
	result := Response{}
	err := OcrTemplate("./files/date-template.png", &result)
	assert.Nil(t, err)
	t.Log(result)
}
func TestOcrGeneral(t *testing.T) {
	result := Response{}
	err := OcrGeneral("./files/general.jpg", &result)
	assert.Nil(t, err)
	t.Log(result)
}

func TestDocumentProcess(t *testing.T) {
	var err error
	result := Response{}
	err = SubmitDocument("./files/cv.pdf", &result)
	assert.Nil(t, err)
	t.Log(result.TaskID)

	t.Log("wait 10 seconds to get the result")
	time.Sleep(time.Duration(10) * time.Second)
	QueryDocumentResult(result.TaskID, &result)
	t.Log(result)
}

func TestContractProcess(t *testing.T) {
	var err error
	result := Response{}
	err = SubmitContract("./files/contract1.docx", "./files/contract2.png", &result)
	assert.Nil(t, err)
	t.Log(result.TaskID)

	taskID := result.TaskID

	t.Log("wait 10 seconds to get the result")
	time.Sleep(time.Duration(10) * time.Second)
	QueryContractResult(taskID, &result)
	t.Log(result)

	t.Log("wait 10 seconds to get the download link")
	DownloadContract(taskID, &result)
	t.Log(result)
}

func TestFlowProcess(t *testing.T) {
	var err error
	result := Response{}
	err = SubmitFlow("./files/cv.pdf", "cv.pdf", &result)
	assert.Nil(t, err)
	t.Log(result.TaskID)

	t.Log("wait 10 seconds to get the result")
	time.Sleep(time.Duration(10) * time.Second)
	QueryFlowResult(result.TaskID, true, &result)
	t.Log(result)
}
