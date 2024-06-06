package getters;

import (
	"os"
	"log"
	v "github.com/DOMIN1310/webmake/vars"
	"context"
)

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