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
	"path"
	"strconv"
	"time"

	v "github.com/DOMIN1310/webmake/vars"
)

func Search(ch chan string, property string, ctx context.Context){
	switch property{
	case "index.ts":
		if req, err := http.NewRequest(
			"GET",
			"https://raw.githubusercontent.com/DOMIN1310/webmake/master/res/tsconfig.json",
			nil,
		); err != nil {
			log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, err.Error());
		} else {
			if res, err := http.DefaultClient.Do(req); err != nil {
				log.Fatalf("%v:%v%v\n", v.ERROR, v.ERROR, err.Error());
			} else {
				if buffer, err := io.ReadAll(res.Body); err != nil {
					log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, err.Error());
				} else {
					if err := os.Mkdir("ts", 0766); err != nil {
						log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, "could not create ts directory");
					} else {
						if err := os.WriteFile(path.Join("tsconfig.json"), buffer, 0744); err != nil{
							log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, "could not create and overwrite ts/tsconfig.json");
						} else {
							if err := os.WriteFile("ts/index.ts", []byte{}, 0744); err != nil{
								log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, err.Error())
							} else {
								if ctx.Err() != nil{
									log.Printf("%v:%v%v\n", v.WARN, v.RESET, "UNABLE TO REACH THAT CODE!");
								} else {
									ch<- property;
								}
							}
						}
					}
				}
			}
		}
	case "tailwindutils.css":
		if _, err := Cmd(func() *exec.Cmd {
			return exec.Command("pnpm", "install", "-D", "tailwindcss")
		}()); err != nil{
			log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, err.Error());
		} else {
			log.Printf("%v:%v%v\n", v.SUCCESS, v.RESET, "successfully initialized tailwindcss");
		}
		if req, err := http.NewRequest(
				"GET", 
				"https://raw.githubusercontent.com/DOMIN1310/webmake/master/res/tailwind.config.js",
				nil,
			); err != nil{
				log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, err.Error());
			} else {
				if res, err := http.DefaultClient.Do(req); err != nil {
					log.Fatalf("%v:%v%v\n", v.ERROR, v.ERROR, err.Error());
				} else {
					if buffer, err := io.ReadAll(res.Body); err != nil {
						log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, err.Error());
					} else {
						if err := os.WriteFile("tailwind.config.js", buffer, 0744); err != nil{
							log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, "could not create and overwrite tailwind.config.js");
						} else {
							if ctx.Err() != nil{
								log.Printf("%v:%v%v\n", v.WARN, v.RESET, "UNABLE TO REACH THAT CODE!");
							} else {
								if err := os.Mkdir("css", 0766); err != nil {
									log.Printf("%v:%v%v\n", v.WARN, v.RESET, "unable to create css directory!");
								} else {
									if err := os.WriteFile(
										"css/tailwind.css",
										[]byte("@tailwind base;\n@tailwind components;\n@tailwind utilities;"),
										0744,
									); err != nil {
										log.Printf("%v:%v%v\n", v.WARN, v.RESET, "unable to create and overwrite css/tailwind.css");
									} else {
										ch<-"tailwind";
									}
								}
							}
						}
					}
				}
			}
	case "main.scss":
		if _, err := Cmd(func() *exec.Cmd{
			return exec.Command("pnpm", "install", "sass", "--save-dev");
		}()); err != nil {
			log.Printf("%v:%v%v\n", v.WARN, v.RESET, err.Error());
		} else {
			log.Printf("%v:%v%v\n", v.SUCCESS, v.RESET, "successfully initialized sass");
			if err := os.Mkdir("sass", 0766); err != nil {
				log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, "could not create sass directory");
			} else {
				if _, err := os.Create("sass/main.scss"); err != nil {
					log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, "error while creating main.scss");
				} else {
					if ctx.Err() != nil{
						log.Printf("%v:%v%v\n", v.WARN, v.RESET, "UNABLE TO REACH THAT CODE!");
					} else {
						ch<-property;
					}
				}
			}
		}
	default:
		if property == "index.html" || property == "index.php"{
			if err := os.Mkdir("public", 0766); err != nil{
				log.Printf("%v:%v%v\n", v.WARN, v.RESET, "unable to create public directory");
			} else {
				log.Printf("%v:%v%v\n", v.SUCCESS, v.RESET, "successfully created public directory");
				var url string = fmt.Sprintf("https://raw.githubusercontent.com/DOMIN1310/webmake/master/res/%v", property);
				if req, err := http.NewRequest("GET", url, nil); err != nil {
					log.Printf("%v:%v%v\n", v.WARN, v.RESET, err.Error());
				} else {
					if res, err := http.DefaultClient.Do(req); err != nil{
						log.Printf("%v:%v%v\n", v.WARN, v.RESET, err.Error());
					} else {
						if buffer, err := io.ReadAll(res.Body); err != nil{
							log.Printf("%v:%v%v\n", v.WARN, v.RESET, err.Error());
						} else {
							if err := os.WriteFile(path.Join("public", property), buffer, 0744); err != nil {
								log.Printf("%v:%v%v\n", v.WARN, v.RESET, "could not create or overwrite " + property);
							} else {
								if ctx.Err() != nil {
									log.Printf("%v:%v%v\n", v.WARN, v.RESET, "UNABLE TO REACH THAT CODE!");
								} else {
									ch <- property;
								}
							}
						}
					}
				}
			}
		} else if property == "index.js" {
			if err := os.Mkdir("js", 0766); err != nil {
				log.Printf("%v:%v%v\n",v.WARN, v.RESET, "could not create js directory")
			} else {
				if _, err := os.Create("js/index.js"); err != nil {
					log.Printf("%v:%v%v\n", v.WARN, v.RESET, "could not create index.js");
				}
			}
		} else if property == "main.css" {
			if err := os.Mkdir("css", 0766); err != nil {
				log.Printf("%v:%v%v\n",v.WARN, v.RESET, "could not create css directory")
			} else {
				if _, err := os.Create("css/main.css"); err != nil {
					log.Printf("%v:%v%v\n", v.WARN, v.RESET, "could not create main.css");
				}
			}
		}
	}
}

func Cmd(f *exec.Cmd) ([]byte, error) {
	if output, err := f.Output(); err != nil {
		return nil, errors.New("incorrect command");
	} else {
		return output, nil;
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
				var chfile chan string = make(chan string, 96);
				var ctx, deadline = context.WithDeadline(context.Background(), time.Now().Add(10*time.Second));
				defer deadline();
				if _, err := Cmd(func() *exec.Cmd{
					return exec.Command("pnpm", "init");
				}()); err != nil {
					log.Printf("%v:%v%v\n", v.WARN, v.RESET, err.Error());
				} else {
					log.Printf("%v:%v%v\n", v.SUCCESS, v.RESET, "successfully initialized package.json");
				}
				go Search(chfile, conf.Findex, ctx);
				go Search(chfile, conf.Styleindex, ctx);
				go Search(chfile, conf.Tmplindex, ctx);
				for {
					select {
					case <- ctx.Done():
						log.Printf("%v:%v%v\n", v.WARN, v.RESET, "contexted has finished");
						return nil;
					case msg := <- chfile:
						log.Printf("%v:%v%v\n", v.SUCCESS, v.RESET, msg);
					}
				}
			}
		}
	}
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
	if style == "css" {
		style = "main.css";
	} else if style == "scss" {
		style = "main.scss";	
	} else if style == "tailwind" {
		style = "tailwindutils.css";
	}
	//marshal data
	var buffer, err = json.Marshal(&v.Template{
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