package payments

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

const (
	ApiGateway      = "https://api.yookassa.ru/v3/payments"
	ApiCheckPayment = "https://api.yookassa.ru/v3/payments/%v"
)

type Yokassa struct {
	ShopId int
	Secret string
	log    *logrus.Logger
}

func NewYokassa(shopId int, secret string, logger *logrus.Logger) Yokassa {
	return Yokassa{
		ShopId: shopId,
		Secret: secret,
		log:    logger,
	}
}

func (y Yokassa) PaymentRequest(yooPayment YokassaPayment, idempotence string) (Payment, error) {
	jsonBody, err := json.Marshal(yooPayment)
	if err != nil {
		y.log.Log(
			logrus.ErrorLevel,
			"Error in payments.Yokassa.PaymentRequest: "+err.Error(),
		)
		return Payment{}, err
	}
	body := bytes.NewBuffer(jsonBody)
	req, err := http.NewRequest("POST", ApiGateway, body)
	if err != nil {
		y.log.Log(
			logrus.ErrorLevel,
			"Error in payments.Yokassa.PaymentRequest: "+err.Error(),
		)
		return Payment{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	shopStr := strconv.Itoa(y.ShopId)
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(shopStr+":"+y.Secret)))
	req.Header.Add("Idempotence-Key", idempotence)
	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		y.log.Log(
			logrus.ErrorLevel,
			"Error in payments.Yokassa.PaymentRequest: "+err.Error(),
		)
		return Payment{}, err
	}
	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		y.log.Log(
			logrus.ErrorLevel,
			"Error in payments.Yokassa.PaymentRequest: "+err.Error(),
		)
		return Payment{}, err
	}
	var payment Payment
	json.Unmarshal(bodyResp, &payment)
	return payment, nil
}

func (y Yokassa) CheckPaymentRequest(uid string) (Payment, error) {
	url := fmt.Sprintf(ApiCheckPayment, uid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		y.log.Log(
			logrus.ErrorLevel,
			"Error in payments.Yokassa.CheckPaymentRequest: "+err.Error(),
		)
		return Payment{}, err
	}
	shopStr := strconv.Itoa(y.ShopId)
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(shopStr+":"+y.Secret)))
	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		y.log.Log(
			logrus.ErrorLevel,
			"Error in payments.Yokassa.CheckPaymentRequest: "+err.Error(),
		)
		return Payment{}, err
	}
	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		y.log.Log(
			logrus.ErrorLevel,
			"Error in payments.Yokassa.CheckPaymentRequest: "+err.Error(),
		)
		return Payment{}, err
	}
	var payment Payment
	json.Unmarshal(bodyResp, &payment)
	return payment, nil
}
