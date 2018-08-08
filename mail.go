package main

import (
	"time"

	"bytes"

	"fmt"

	"github.com/mytokenio/go_sdk/log"
	"gopkg.in/gomail.v2"
)

var (
	MailClient  *gomail.Dialer
	mailContent = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
	<html xmlns="http://www.w3.org/1999/xhtml">
	<head>
	<style>
        table,table tr th, table tr td { border:1px solid #0094ff;padding:2px 15px; }
        table {min-height: 25px; line-height: 25px; text-align: center; border-collapse: collapse;}   
    </style>
	</head>
	<body>
		<table>
				<thead>
			        <tr>
			          <th>时间</th>
			          <th>买入均价</th>
			          <th>卖出均价</th>
			          <th>差价</th>
			          <th>盈利</th>
			        </tr>
		      	</thead>
		      	<tbody>
					%s
		      	</tbody>
		    </table>
	</body>
	</html>
`
	mailTable = "<tr style='%s'><td>%s</td><td>%.2f</td><td>%.2f</td><td>%.2f</td><td>%.2f%%</td></tr>"
)

func sendMail(list []*A) {
	content := initMessage(list)
	MailClient = gomail.NewDialer("smtp.exmail.qq.com", 465, "me@wen.moe", "RMzsBXFYMmw79j")
	m := gomail.NewMessage()
	m.SetHeader("From", "me@wen.moe")
	m.SetHeader("To", "me@wen.moe", "li.1328@osu.edu")
	//m.SetHeader("To", "me@wen.moe")

	m.SetHeader("Subject", time.Now().Add(-time.Hour*24).Format("2006-01-02")+"每日行情")
	m.SetBody("text/html", fmt.Sprintf(mailContent, content))
	if err := MailClient.DialAndSend(m); err != nil {
		log.Error(" MailClient.DialAndSend() err(%v)", err)
	}
}

func initMessage(list []*A) string {
	var (
		bb        bytes.Buffer
		isSuccess string
	)

	for _, l := range list {
		if ((l.SellPer-l.BuyPer)*l.Vol)/totalPrice > 0.04 {
			isSuccess = "color:green"
		}
		ss := fmt.Sprintf(mailTable, isSuccess, l.TS.Format("2006-01-02  15:04"), l.BuyPer, l.SellPer, l.SellAll-totalPrice, (l.SellAll-totalPrice)/totalPrice*100)
		bb.WriteString(ss)
	}
	return bb.String()
}
