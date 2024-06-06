package getters;

import (
	"context"
	"net/http"
	"log"
	v "github.com/DOMIN1310/webmake/vars"
	"io"
	"encoding/json"
)

func GetScripts(ctx context.Context) v.Scripts {
	if req, err := http.NewRequestWithContext(
		ctx, 
		"GET", 
		"https://raw.githubusercontent.com/DOMIN1310/webmake/master/res/scripts.json",
		nil,
	); err != nil {
		log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, "unable to get scripts request");
		return v.Scripts{};		
	} else {
		if res, err := http.DefaultClient.Do(req); err != nil {
			log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, "unable to get response from the request to get scripts");
			return v.Scripts{};
		} else {
			if buffer, err := io.ReadAll(res.Body); err != nil {
				res.Body.Close();
				log.Printf("%v:%v%v\n", v.ERROR, v.RESET, "unable to read response");
				return v.Scripts{};
			} else {
				defer res.Body.Close();
				var variable v.Scripts;
				if err := json.Unmarshal(buffer, &variable); err != nil {
					log.Printf("%v:%v%v\n", v.ERROR, v.RESET, "unable to marshal data");
					return v.Scripts{};
				} else {
					return variable;
				}
			}
		}
	}
}