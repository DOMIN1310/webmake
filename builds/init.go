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
	"strconv"
	"strings"
	"time"

	g "github.com/DOMIN1310/webmake/getters"
	v "github.com/DOMIN1310/webmake/vars"
)

func Search(ch chan string, property string, ctx context.Context){
	switch property{
	case "ts/index.ts":
		var url string = "https://raw.githubusercontent.com/DOMIN1310/webmake/master/res/tsconfig.json";
		g.PrepareComponent(ch, "ts", map[string]string{property: ""}, ctx);
		g.PrepareComponent(ch, "./", g.GetURLBody(url, ctx, "GET", "tsconfig.json"), ctx);
	case "css/tailwindutils.css":
		if err := g.Cmd(func() *exec.Cmd {
			return exec.Command("pnpm", "install", "-D", "tailwindcss")
		}()); err != nil{
			log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, err.Error());
		} else {
			log.Printf("%v:%v%v\n", v.INIT, v.RESET, "successfully initialized tailwindcss");
		}
		var url string = "https://raw.githubusercontent.com/DOMIN1310/webmake/master/res/tailwind.config.js";
		g.PrepareComponent(ch, "./", g.GetURLBody(url, ctx, "GET", "tailwind.config.js"), ctx);
		g.PrepareComponent(ch, "css", map[string]string{
			property: "@tailwind base;\n@tailwind components;\n@tailwind utilities;",
		}, ctx);
	case "sass/main.scss":
		if err := g.Cmd(func() *exec.Cmd{
			return exec.Command("pnpm", "install", "sass", "--save-dev");
		}()); err != nil {
			log.Printf("%v:%v%v\n", v.WARN, v.RESET, err.Error());
		} else {
			log.Printf("%v:%v%v\n", v.INIT, v.RESET, "successfully initialized sass");
			g.PrepareComponent(ch, "sass", map[string]string{
				property: "",
			}, ctx);
		}
	default:
		g.PrepareComponent(ch, func () string {
			if property == "css/main.css" {
				return "css";
			} else if property == "js/index.js" {
				return "js";
			} else {
				return "public";
			}
		}(), func () map[string]string {
			if property == "public/index.html" || property == "public/index.php" {
				var url string = fmt.Sprintf("https://raw.githubusercontent.com/DOMIN1310/webmake/master/res/%v", strings.Split(property, "/")[1]);
				return g.GetURLBody(url, ctx, "GET", property);
			} else {
				return map[string]string{property: ""};
			}
		}(), ctx);
	}
}

func createWeb(ctx context.Context) error {
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
				return errors.New("unmarshalation error");
			} else {
				var chfile chan string = make(chan string, 96);
				if err := g.Cmd(func() *exec.Cmd{
					return exec.Command("pnpm", "init");
				}()); err != nil {
					log.Printf("%v:%v%v\n", v.WARN, v.RESET, err.Error());
				} else {
					log.Printf("%v:%v%v\n", v.INIT, v.RESET, "successfully initialized package.json");
				}
				if conf.Git {
					if err := g.Cmd(func () *exec.Cmd {
						return exec.Command("git", "init");
					}()); err != nil {
						log.Printf("%v:%v%v\n", v.ERROR, v.RESET, "unable to initialize git ensure u have git installed");
					} else {
						log.Printf("%v:%v%v\n", v.INIT, v.RESET,"successfully initialized git!")
					}
				}
				go Search(chfile, conf.Findex, ctx);
				go Search(chfile, conf.Styleindex, ctx);
				go Search(chfile, conf.Tmplindex, ctx);
				for i := 0; i <= 3; i++{
					select {
					case <- ctx.Done():
						log.Printf("%v:%v%v\n", v.FINISHED, v.RESET, "context has finished");
						return nil;
					case msg := <- chfile:
						log.Printf("%v:%v%v\n", v.SUCCESS, v.RESET, fmt.Sprintf("successfully created %v!", msg));
					}
				}
			}
		}
	}
	return nil;
}

func InitPackage() {
	//get inputs
		//vars
	var flang string;
	var style string;
	var tmpl string;
	var temp string;
	var git bool;
	//scan them
	fmt.Println("functional language [js, ts]: ");
		//flang
	fmt.Scan(&flang);
	if flang != "ts" && flang != "js" {
		log.Fatalf("%v:%v%v", v.ERROR, v.RESET, "invalid option choose either js or ts");
	}
	fmt.Println("stylesheet [basic (css), scss, tailwind]: ");
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
	if style == "basic" {
		style = "css/main.css";
	} else if style == "scss" {
		style = "sass/main.scss";	
	} else if style == "tailwind" {
		style = "css/tailwindutils.css";
	}
	if flang == "ts" {
		flang = "ts/index.ts"
	} else if flang == "js" {
		flang = "js/index.js"
	}
	var ctx, deadline = context.WithDeadline(context.Background(), time.Now().Add(10*time.Second));
	defer deadline();
	scripts := g.GetScripts(ctx);
	if buffer, err := json.MarshalIndent(&v.Template{
		Findex: flang,
		Styleindex: style,
		Tmplindex: "public/index." + tmpl,
		Git: git,
		Scripts: scripts,
	}, "", "  "); err != nil {
		log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, "unable to marshal wb-package.json with scripts!");
	} else {
		log.Printf("%v:%v%v\n", v.SUCCESS, v.RESET, "Marshaling was successful!");
		g.PrepareComponent(nil, "./", map[string]string{"wb-package.json": string(buffer)}, nil)
		if err := createWeb(ctx); err != nil {
			log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, err.Error());
		} else {
			log.Printf("%v:%v%v\n", v.DONE, v.RESET, "INITIALIZATION DONE!");
		}
	}
}