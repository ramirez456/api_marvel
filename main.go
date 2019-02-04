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
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Creation    string  `json:"modified"`
	Serie       Series  `json:"series"`
	Comic       Comics  `json:"comics"`
	Storie      Stories `json:"stories"`
}
type Comics struct {
	Items []Item `json:"items"`
}
type Series struct {
	Items []Item_s `json:"items"`
}
type Stories struct {
	Items []Item `json:"items"`
}
type Item struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
type Item_s struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func main() {
	//Timestamp
	t := time.Now()
	ts := t.Format("20060102150405")
	//Creamos el hash
	hasher := md5.New()
	hasher.Write([]byte(ts + api_priv + api_key))
	//Pasamos el hash a texto
	texto := hex.EncodeToString(hasher.Sum(nil))
	//Creamos el menu
	menu :=
		`
	Bienvenido
	
	¿Qué deseas hacer?

	[ 1 ] Personaje
	[ 2 ] Listarlos
	`
	fmt.Print(menu)

	option := 0

	fmt.Scanf("%d", &option)

	switch option {

	case 1:

		fmt.Println("¿Qué pensonajes deseas ver?")

		reader := bufio.NewReader(os.Stdin)
		name, _ := reader.ReadString('\n')

		fmt.Print("Los datos de " + name)
		//pasamos el nombre a modo URL " "="%20"
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
		fmt.Println(objeto.Data.Heros[0].Id)
		fmt.Println("*--Nombre--*")
		fmt.Println(objeto.Data.Heros[0].Name)
		fmt.Println("*--Description--*")
		fmt.Println(objeto.Data.Heros[0].Description)
		fmt.Println("*--Fecha de Creación--*")
		fmt.Println(objeto.Data.Heros[0].Creation)
		fmt.Println("*--Comics--*")
		cant := len(objeto.Data.Heros[0].Comic.Items)
		fmt.Print("*--Existen ")
		fmt.Print(cant)
		fmt.Println(" Comics-*")

		for i := 0; i < cant; i++ {
			fmt.Println(objeto.Data.Heros[0].Comic.Items[i].Name)
		}
		fmt.Println("*--Stories--*")
		canti := len(objeto.Data.Heros[0].Storie.Items)
		fmt.Print("*--Existen ")
		fmt.Print(canti)
		fmt.Println(" Stories-*")
		for i := 0; i < canti; i++ {
			fmt.Println(objeto.Data.Heros[0].Storie.Items[i].Name)
		}
		fmt.Println("*--Series--*")
		cantid := len(objeto.Data.Heros[0].Serie.Items)
		fmt.Print("*--Existen ")
		fmt.Print(cantid)
		fmt.Println(" Series-*")
		for i := 0; i < cantid; i++ {
			fmt.Println(objeto.Data.Heros[0].Serie.Items[i].Name)
		}

	case 2:

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
