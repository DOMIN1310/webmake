package builds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	v "github.com/DOMIN1310/webmake/vars"
)

func Search(ch chan string, property string, ctx context.Context){
	switch property{
	case "ts/index.ts":
		var url string = "https://raw.githubusercontent.com/DOMIN1310/webmake/master/res/tsconfig.json";
		PrepareComponent(ch, "ts", map[string]string{property: ""}, ctx);
		PrepareComponent(ch, "./", GetURLBody(url, ctx, "GET", property), ctx);
	case "css/tailwindutils.css":
		if err := Cmd(func() *exec.Cmd {
			return exec.Command("pnpm", "install", "-D", "tailwindcss")
		}()); err != nil{
			log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, err.Error());
		} else {
			log.Printf("%v:%v%v\n", v.INIT, v.RESET, "successfully initialized tailwindcss");
		}
		var url string = "https://raw.githubusercontent.com/DOMIN1310/webmake/master/res/tailwind.config.js";
		PrepareComponent(ch, "./", GetURLBody(url, ctx, "GET", "tailwind.config.js"), ctx);
		PrepareComponent(ch, "css", map[string]string{
			property: "@tailwind base;\n@tailwind components;\n@tailwind utilities;",
		}, ctx);
	case "sass/main.scss":
		if err := Cmd(func() *exec.Cmd{
			return exec.Command("pnpm", "install", "sass", "--save-dev");
		}()); err != nil {
			log.Printf("%v:%v%v\n", v.WARN, v.RESET, err.Error());
		} else {
			log.Printf("%v:%v%v\n", v.INIT, v.RESET, "successfully initialized sass");
			PrepareComponent(ch, "sass", map[string]string{
				property: "",
			}, ctx);
		}
	default:
		PrepareComponent(ch, func () string {
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
				return GetURLBody(url, ctx, "GET", property);
			} else {
				return map[string]string{property: ""};
			}
		}(), ctx);
	}
}

func GetURLBody(url string, ctx context.Context, method string, file string) map[string]string {
	if req, err := http.NewRequestWithContext(ctx, method, url, nil); err != nil {
		log.Printf("%v:%v%v\n", v.WARN, v.RESET, "unable to get any request within context and destiny url");
		return map[string]string{file: ""};	
	} else if ctx.Err() != nil {
		fmt.Printf("%v:%v%v\n", v.WARN, v.RESET, "unable to complete request: context error")
		return map[string]string{file: ""};
	}else {
		if res, err := http.DefaultClient.Do(req); err != nil {
			log.Printf("%v:%v%v\n", v.WARN, v.RESET, "unable to complete request with proper response");
			return map[string]string{file: ""};
		} else {
			if buffer, err := io.ReadAll(res.Body); err != nil {
				log.Printf("%v:%v%v\n", v.WARN, v.RESET, "unable to read response body!");
				return map[string]string{file: ""};
			} else {
				return map[string]string{
					file: string(buffer),
				}
			}						
		}
	}
}

func PrepareComponent(delch chan string, dirName string, files map[string]string, ctx context.Context) {
	if dirName != "./" {
    if _, err := os.Stat(dirName); err != nil {
			if os.IsNotExist(err) {
				if err := os.Mkdir(dirName, 0755); err != nil {
					log.Printf("%v:%v%v\n", v.WARN, v.RESET, "unable to create " + dirName + " directory")
				} else {
					log.Printf("%v:%v%v\n", v.CREATION, v.RESET, "successfully created " + dirName);
				}
			} else {
				log.Printf("%v:%v%v\n", v.ERROR, v.RESET, "unable to check if directory exists or not!")
				return
			}
    } else if ctx.Err() != nil {
			log.Printf("%v:%v%v", v.ERROR, v.RESET, "context error!! unable to reach the code further!")
    }
	}
	for fileName, content := range files {
		if err := os.WriteFile(fileName, []byte(content), 0755); err != nil {
			log.Printf("%v:%v%v\n", v.WARN, v.RESET, "could not create " + fileName);
		} else {
			if delch != nil {
				delch<-fileName;
			}
		}
	}
}

func Cmd(f *exec.Cmd) error {
	if err := f.Run(); err != nil {
		return errors.New("incorrect command");
	} else {
		return nil;
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
				return errors.New("unmarshalation error");
			} else {
				var chfile chan string = make(chan string, 96);
				var ctx, deadline = context.WithDeadline(context.Background(), time.Now().Add(10*time.Second));
				defer deadline();
				if err := Cmd(func() *exec.Cmd{
					return exec.Command("pnpm", "init");
				}()); err != nil {
					log.Printf("%v:%v%v\n", v.WARN, v.RESET, err.Error());
				} else {
					log.Printf("%v:%v%v\n", v.INIT, v.RESET, "successfully initialized package.json");
				}
				if conf.Git {
					if err := Cmd(func () *exec.Cmd {
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
	//marshal data
	var buffer, err = json.MarshalIndent(&v.Template{
		Findex: flang,
		Styleindex: style,
		Tmplindex: "public/index." + tmpl,
		Git: git,
	}, "", "  ");
	//check if marshaling was successful
	if err != nil{
		log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, err.Error());
	} else {
		log.Printf("%v:%v%v\n", v.SUCCESS, v.RESET, "Marshaling was successful!");
	}
	PrepareComponent(nil, "./", map[string]string{"wb-package.json": string(buffer)}, nil)
	if err := createWeb(); err != nil {
		log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, "UNABLE TO INITILIAZE THE PROJECT")
	} else {
		log.Printf("%v:%v%v\n", v.DONE, v.RESET, "INITIALIZATION DONE!");
	}
}