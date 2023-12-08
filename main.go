package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

type API_data struct {
	URL string `json:"url"`
}

type Conf struct {
	API_key string 
}

func main() {
	config , err := Get_config()
	if( err != nil ) {
		log.Fatal(err)
	}

	url := "https://api.nasa.gov/planetary/apod?api_key=" + config.API_key

	client := http.Client{}
	resp , err := client.Get(url)
	if( err != nil ) {
		log.Fatal(err)
	} 
	defer resp.Body.Close()

	ansewr , err := io.ReadAll(resp.Body)
	if( err != nil ) {
		log.Fatal(err)
	}

	if( resp.StatusCode != 200 ) {
		fmt.Printf( "status code:%d\n%s" , resp.StatusCode , ansewr)
		os.Exit(1)
	}

	data := API_data{}
	err = json.Unmarshal( ansewr , &data )
	if( err != nil ) {
		log.Fatal(err)
	}

	err = Download_picture( &client , &data )
	if( err != nil ) {
		log.Fatal(err)
	}

	err = Set_background()
	if( err != nil ) {
		log.Fatal(err)
	}
}

func Get_config() (Conf , error) {
	handle := viper.New()
	handle.SetConfigType("yaml")
	handle.AddConfigPath( "." )
	handle.SetConfigName( "config.yaml" )
	
	err := handle.ReadInConfig()
	if( err != nil ) {
		return Conf{} , err
	}

	conf := Conf{}
	err = handle.Unmarshal(&conf)
	if( err != nil ) {
		return Conf{} , err
	}

	return conf , nil
}

func Download_picture( client *http.Client , data *API_data ) error {
	buffer , err := client.Get(data.URL)
	if( err != nil ) {
		return err
	}

	file , err := os.Create( "apod.jpg" )
	if( err != nil ) {
		return err
	}

	_ , err = io.Copy( file , buffer.Body )
	if( err != nil ) {
		return err
	}

	return nil
}

func Set_background() error {
	dir , err := os.Getwd()
	if( err != nil ) {
		return err
	}

	file_name := fmt.Sprintf( "\"file://%s/apod.jpg\"" , dir )
	fmt.Print(file_name)
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", file_name )
	err = cmd.Run()
	if( err != nil ) {
		return err
	}

	return nil
}