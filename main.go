package main //сделать структуру(CurrencyConvertes Convert()), с методом, покрыть тестами // https://pkg.go.dev/net/http/httptest //докер

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/text/encoding/charmap"
)

const (
	currenciesURL = "http://www.cbr12343221.ru/scripts/XML_daily.asp"
	cur           = `JPY`
)

type Currency struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

type ValCurs struct {
	XMLName     xml.Name   `xml:"ValCurs"`
	XMLDate     string     `xml:"Date,attr"`
	XMLNameAttr string     `xml:"name,attr"`
	Currencies  []Currency `xml:"Valute"`
}

func getBody(url string) ([]byte, error) { //todo сделать таймаут, вынести в константу

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func decodeWindows1251(date []byte) ([]byte, error) {
	dec := charmap.Windows1251.NewDecoder()
	out, err := dec.Bytes(date)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func main() {

	date, err := getBody(currenciesURL)

	if err != nil {
		log.Fatal(err)
	}

	valCurs := &ValCurs{}

	temp, err := decodeWindows1251(date[1:]) // пропускаем символ '<', чтобы не учитывался тег общей инфы (о кодировке, версии xml)
	if err != nil {
		log.Fatal(err)
	}

	if err := xml.Unmarshal(temp, valCurs); err != nil {
		log.Fatal(err)
	}

	for i := len(valCurs.Currencies) - 1; i >= 0; i-- {
		if valCurs.Currencies[i].CharCode == cur {
			fmt.Printf("Курс %s %s к российскому рублю составляет %s \nна %s\n",
				valCurs.Currencies[i].Nominal,
				valCurs.Currencies[i].Name,
				valCurs.Currencies[i].Value,
				valCurs.XMLDate)
			break
		}
	}

}
