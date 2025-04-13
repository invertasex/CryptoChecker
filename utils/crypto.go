package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var symbolMap = map[string]string{
	"BTC":  "bitcoin",
	"ETH":  "ethereum",
	"BNB":  "binancecoin",
	"USDT": "tether",
	"TRX":  "tron",
	"SOL":  "solana",
}

func GetPrice(symbol string) (float64, error) {
	coinID, ok := symbolMap[symbol]
	if !ok {
		return 0, fmt.Errorf("неизвестный тикер: %s", symbol)
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", coinID)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data map[string]map[string]float64
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, err
	}

	price := data[coinID]["usd"]
	return price, nil
}

func GetTRC20Balance(wallet string) (float64, error) {
	if !strings.HasPrefix(wallet, "T") || len(wallet) < 34 {
		return 0, fmt.Errorf("некорректный TRON-адрес")
	}

	url := fmt.Sprintf("https://apilist.tronscan.org/api/account?address=%s", wallet)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	rawTokens, ok := result["trc20token_balances"]
	if !ok {
		return 0, fmt.Errorf("на этом адресе не найдено TRC-20 токенов")
	}

	tokens, ok := rawTokens.([]interface{})
	if !ok {
		return 0, fmt.Errorf("неожиданный формат токенов")
	}

	for _, t := range tokens {
		tok, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		if tok["tokenName"] == "Tether USD" {
			balanceStr, ok := tok["balance"].(string)
			if !ok {
				continue
			}
			balance, _ := strconv.ParseFloat(balanceStr, 64)
			return balance / 1e6, nil
		}
	}

	return 0, fmt.Errorf("на кошельке нет USDT")
}
