package main

import (
	"flag"
	"log"
	b "github.com/DOMIN1310/webmake/builds"
)

func main(){
	var init *bool = flag.Bool("init", false, "[true, false]");
	var run *string = flag.String("run", "", "example: ./webmake -run=category:script look for scripts in wb-package.json")
	flag.Parse();
	if *run != "" && *init{
		log.Fatalf("too many flags");
	} else if *init {
		log.Println("initializing the project");
		b.InitPackage();
	} else if *run != "" {
		b.RunScript();
	}
}