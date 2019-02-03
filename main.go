package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const api_key = "2b5ca51c7262be5c2ac6809bc7f13187"
const api_priv = "128fbc552a93125488096774d4b774b4caada843"

type Response struct {
	Total int     `json:"total"`
	Data  Results `json:"data"`
}
type Results struct {
	Heros []Hero `json:"results"`
}
type Hero struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Comic       Comics `json:"comics"`
}
type Comics struct {
	Items []Item `json:"items"`
}
type Item struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func main() {

	t := time.Now()
	ts := t.Format("20060102150405")

	hasher := md5.New()
	hasher.Write([]byte(ts + api_priv + api_key))
	texto := hex.EncodeToString(hasher.Sum(nil))

	link := "http://gateway.marvel.com/v1/public/characters?ts=" + ts + "&apikey=" + api_key + "&hash=" + texto

	response, err := http.Get(link)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	menu :=
		`
	Bienvenido
	
	Bienvenido
	Que deseas hacer
	[ 1 ] Personaje
	[ 2 ] Listarlos
	`
	fmt.Print(menu)

	option := 0
	fmt.Scanf("%d", &option)
	switch option {
	case 1:
		fmt.Println("Que pensonajes deseas ver")

		reader := bufio.NewReader(os.Stdin)
		name, _ := reader.ReadString('\n')

		fmt.Print("Los datos de " + name)

		uname := &url.URL{}
		err := uname.UnmarshalBinary([]byte(name))
		if err != nil {
			log.Fatal(err)
		}

		nombre := uname.String()
		nom := strings.Replace(nombre, "%0A", "", -1)

		uri := "http://gateway.marvel.com/v1/public/characters?name=" + nom + "&ts=" + ts + "&apikey=" + api_key + "&hash=" + texto

		resp, err := http.Get(uri)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		respData, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		var objeto Response
		json.Unmarshal(respData, &objeto)

		fmt.Println("*--ID--*")
		fmt.Println(responseObject.Data.Heros[0].Id)
		fmt.Println("*--Nombre--*")
		fmt.Println(responseObject.Data.Heros[0].Name)
		fmt.Println("*--Description--*")
		fmt.Println(responseObject.Data.Heros[0].Description)
		fmt.Println("*--Comics--*")
		cant := len(responseObject.Data.Heros[0].Comic.Items)

		for i := 0; i < cant; i++ {
			fmt.Println(responseObject.Data.Heros[0].Comic.Items[i].Name)
		}

	case 2:
		fmt.Println("Los personajes de marvel son:")
		for i := 0; i < len(responseObject.Data.Heros); i++ {
			fmt.Print("ID: ")
			fmt.Print(responseObject.Data.Heros[i].Id)
			fmt.Print(" Nombre: " + responseObject.Data.Heros[i].Name)
			fmt.Println(" *--Descripción: " + responseObject.Data.Heros[i].Description)
		}
	default:
		fmt.Println("Esta opción no es validad tienes precionar 1 ó 2")
	}

}
