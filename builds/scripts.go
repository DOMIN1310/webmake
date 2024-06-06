package builds

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	g "github.com/DOMIN1310/webmake/getters"
	v "github.com/DOMIN1310/webmake/vars"
)

func Readwebmakepackage(ctx context.Context) g.Scripts {
	if file, err := os.Open("./wb-package.json"); err != nil {
		file.Close();
		log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, "unable to open wb-package.json ensure you are in the same directory as the file is");
		return g.Scripts{};
	} else {
		defer file.Close();
		if buffer, err := io.ReadAll(file); err != nil {
			log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, "unable to react wb-package.json");
			return g.Scripts{};
		} else {
			var tmpl v.Template;
			if err := json.Unmarshal(buffer, &tmpl); err != nil {
				log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, "unable to unmarshal buffer");
				return g.Scripts{};
			} else {
				return tmpl.Scripts;
			}
		}
	}
}

func RunScript(ctx context.Context, arg string) {
	var scripts g.Scripts = Readwebmakepackage(ctx);
	log.Println(scripts);
	log.Println(arg);
}