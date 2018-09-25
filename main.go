package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	dayList    []*A
	totalPrice float64 = 5000
	kraUrl             = "https://api.kraken.com/0/public/Depth?pair=XXBTZCAD"
	quaUrl             = "https://api.quadrigacx.com/v2/order_book?book=btc_cad"
)

type A struct {
	TS      time.Time
	BuyPer  float64 // 买入均价
	SellPer float64 // 卖出均价
	SellAll float64 // 卖出总价
	Vol     float64
}

func main() {
	var (
		lastDay = time.Now().Day()
		err     error
		kk      = new(KraRet)
	)

	for {
		if lastDay != time.Now().Day() {
			log.Printf("start send mail. data count:%d", len(dayList))
			sendMail(dayList)
			dayList = nil
			lastDay = time.Now().Day()
		}
		a := new(A)
		a.TS = time.Now()
		time.Sleep(time.Second * 3)
		kk, err = reqKra()
		if err != nil || len(kk.Error) > 0 || len(kk.Result.XXBTZCAD.Asks) <= 0 {
			log.Printf("reqKra() err(%v)", err)
			continue
		}
		a.Vol = countKra(kk.Result.XXBTZCAD.Asks)
		a.BuyPer = totalPrice / a.Vol
		dd, err := reqQuad()
		if err != nil || len(dd.Bids) <= 0 {
			log.Printf("reqQuad() err(%v)", err)
			continue
		}
		a.SellAll = countQuad(dd, a.Vol)
		a.SellPer = a.SellAll / a.Vol
		dayList = append(dayList, a)
		if len(dayList) > 2 {
			sendMail(dayList)
		}
		//todo 报警邮件
	}
}

func reqKra() (ret *KraRet, err error) {
	ret = new(KraRet)
	res, err := http.Get(kraUrl)
	if err != nil {
		log.Printf("reqKra http.Get() err(%v)", err)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("reqKra ioutil.ReadAll() err(%v)", err)
		return
	}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		log.Printf("reqKra json.Unmarshal(%s) err(%v)", string(body), err)
		return
	}
	return
}

func countKra(dp [][]interface{}) float64 {
	var (
		v, p           float64
		err            error
		usedPrice      float64
		myVol, leftVol float64
	)
	for _, d := range dp {
		p, err = strconv.ParseFloat(d[0].(string), 64)
		if err != nil {
			log.Printf("strconv.ParseFloat(%s) err(%v)", d[0].(string), err)
			continue
		}
		v, err = strconv.ParseFloat(d[1].(string), 64)
		if err != nil {
			log.Printf("strconv.ParseFloat(%s) err(%v)", d[0].(string), err)
			continue
		}
		price := p * v
		if usedPrice+price >= totalPrice {
			leftVol = (totalPrice - usedPrice) / p
			usedPrice += p * leftVol
			myVol += leftVol
			break
		} else {
			myVol += v
			usedPrice += price
		}
	}
	log.Printf("use %.2f by %.2f coin\n", usedPrice, myVol)
	return myVol
}

func reqQuad() (ret *QuadRet, err error) {
	ret = new(QuadRet)
	res, err := http.Get(quaUrl)
	if err != nil {
		log.Printf("reqQuad http.Get() err(%v)", err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("reqQuad ioutil.ReadAll() err(%v)", err)
		return
	}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		log.Printf("reqQuad json.Unmarshal(%s) err(%v)", string(body), err)
		return
	}
	return
}

func countQuad(dd *QuadRet, myVol float64) float64 {
	var (
		v, p     float64
		err      error
		getPrice float64
	)

	for _, d := range dd.Bids {
		p, err = strconv.ParseFloat(d[0], 64)
		if err != nil {
			log.Printf("strconv.ParseFloat(%s) err(%v)", d[0], err)
			continue
		}
		v, err = strconv.ParseFloat(d[1], 64)
		if err != nil {
			log.Printf("strconv.ParseFloat(%s) err(%v)", d[0], err)
			continue
		}
		log.Printf("left %.2f coin get %.2f\n", myVol, getPrice)
		if myVol-v <= 0 {
			getPrice += p * myVol
			myVol = 0
			log.Printf("left %.2f coin get %.2f \n", myVol, getPrice)
			break
		} else {
			getPrice += p * v
			myVol -= v
		}
	}
	return getPrice
}
