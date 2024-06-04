package builds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	v "github.com/DOMIN1310/webmake/vars"
)

func Search(ch chan string, property string, dir string){
	switch property{
	case "index.ts":
		
	case "tailwindutils.css":
		exec.Command("/bin/bash", "res/runners/tailwindinstall.sh")
	case "main.scss":

	default:
		_, err := os.Create(path.Join(dir, property));
		if err != nil{
			log.Printf("%v:%v%v\n", v.WARN, v.RESET, "Unable to create " + property);
		}
		var buffer []byte;
		if property == "html" || property == "php" {
			file, err := os.Open(path.Join("res", property));
			if err != nil {
				log.Printf("%v:%v%v\n", v.WARN, v.RESET, "unable to open pre made template index");
			} else {
				defer file.Close();
				buffer, err = io.ReadAll(file)
				if err != nil {
					log.Printf("%v:%v%v\n", v.WARN, v.RESET, "unable to read pre made template index");
				} else {
					n, err := file.Write(buffer);
					if err != nil{
						log.Printf("%v:%v%v\n", v.WARN, v.RESET, "unable to write pre made template index");
					} else {
						log.Printf("%v:%v%v\n", v.SUCCESS, v.RESET, fmt.Sprintf("successfully overwritten file with %d bytes", n));
					}
				}
			}
		}
	}
}

func createWeb() error{
	if confFile, err := os.Open("wb-package.json"); err != nil{
		return errors.New("could not open wb-package.json, ensure it exists");	
	} else {
		defer confFile.Close();
		confBuffer, err := io.ReadAll(confFile);
		if err != nil{
			return errors.New("unable to read wb-package.json");	
		} else {
			var conf v.Template;
			if err := json.Unmarshal(confBuffer, &conf); err != nil {
				return errors.New("marshalation error");
			} else {
				var chfile chan string;
				var ctx, deadline = context.WithCancel(context.Background())
				defer deadline();
				err := os.Mkdir(conf.Dir, 0755);
				if err != nil {
					log.Fatalf("%v:%v%v", v.ERROR, v.RESET, "error while creating the app directory!");
				}
				
				select {
				case <- ctx.Done():  
					log.Printf("%v:%v%v\n", v.WARN, v.RESET, "context finished");
				case file := <- chfile:
					log.Printf("%v:%v:%v%v\n", v.SUCCESS, v.RESET, "SUCCESSFULLY CREATED ", file);
				}
			}
		}
	}
	return nil;
}

func InitPackage() {
	//get inputs
		//vars
	var dir string;
	var flang string;
	var style string;
	var tmpl string;
	var temp string;
	var git bool;
	//scan them
		//clear terminal
	exec.Command("/bin/clear");
	fmt.Println("functional language [js, ts]: ");
		//flang
	fmt.Scan(&flang);
	if flang != "ts" && flang != "js" {
		log.Fatalf("%v:%v%v", v.ERROR, v.RESET, "invalid option choose either js or ts");
	}
	fmt.Println("stylesheet [base (css), scss, tailwind]: ");
		//styles
	fmt.Scan(&style);
	if style != "basic" && style != "scss" && style != "tailwind"{
		log.Fatalf("%v:%v%v", v.ERROR, v.RESET, "invalid option choose either basic (css), scss, tailwind");
	}
	fmt.Println("template [html, php]: ");
		//template
	fmt.Scan(&tmpl);
	if tmpl != "html" && tmpl != "php"{
		log.Fatalf("%v:%v%v", v.ERROR, v.RESET, "invalid option choose either html or php");
	}
	fmt.Println("initialize git [true, false]");
		//git
	fmt.Scan(&temp)
	if b, e := strconv.ParseBool(temp); e != nil {
		git = b;
	}
	fmt.Println("app directory name: ");
	// app directory
	fmt.Scan(&dir)
	if dir == "" {
		dir = "app"
	}
	if style == "css" {
		style = "main.css";
	} else if style == "scss" {
		style = "main.scss";	
	} else if style == "tailwind" {
		style = "tailwindutils.css";
	}
	//marshal data
	var buffer, err = json.Marshal(&v.Template{
		Dir: dir,
		Findex: "index." + flang,
		Styleindex: style,
		Tmplindex: "index." + tmpl,
		Git: git,
	});
	//check if marshaling was successful
	if err != nil{
		log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, err.Error());
	} else {
		log.Printf("%v:%v%v\n", v.SUCCESS, v.RESET, "Marshaling was successful!");
	}
	//check if file creation was successful
	if _, err = os.Create("wb-package.json"); err != nil{
		log.Printf("%v:%v%v\n", v.WARN, v.RESET, err.Error());
	} else {
		log.Printf("%v:%v%v\n", v.SUCCESS, v.RESET, "File created successfully!")
	}
	//check if file was sucessfully overwritten 
	if err = os.WriteFile("wb-package.json", buffer, 0666); err != nil{
		log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, err.Error());
	} else {
		log.Printf("%v:%v%v\n", v.SUCCESS, v.RESET, "wb-package overwritten successfully!")
	}
	if err = createWeb(); err != nil{
		log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, err.Error());
	} else {
		log.Printf("%v:%v%v\n", v.SUCCESS, v.RESET, "CREATION COMPLETED SUCCESSFULLY!")
	}
}