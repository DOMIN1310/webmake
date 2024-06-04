package main

import (
	"flag"
	"log"
	b "webmake/builds"
)

func main(){
	var compile *bool = flag.Bool("compile", false, "[true, false]");
	var init *bool = flag.Bool("init", false, "[true, false]");
	flag.Parse();
	if *compile && *init{
		log.Fatalf("too many flags");
	} else if *compile {
		log.Println("compiling");
	} else if *init {
		log.Println("initializing the project");
		b.InitPackage();
	}
}