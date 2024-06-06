package main

import (
	"context"
	"flag"
	"log"
	"time"

	b "github.com/DOMIN1310/webmake/builds"
)

func main(){
	//flags
	var init *bool = flag.Bool("init", false, "[true, false]");
	var run *string = flag.String("run", "", "example: ./webmake -run=category:script look for scripts in wb-package.json")
	flag.Parse();
	//context
	var ctx, deadline = context.WithDeadline(context.Background(), time.Now().Add(7*time.Second));
	defer deadline();
	//run
	if *run != "" && *init{
		log.Fatalf("too many flags");
	} else if *init {
		log.Println("initializing the project");
		b.InitPackage(ctx);
	} else if *run != "" {
		b.RunScript(ctx, *run);
	}
}