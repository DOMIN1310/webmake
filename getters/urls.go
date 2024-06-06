package getters;

import (
	"context"
	"net/http"
	"log"
	v "github.com/DOMIN1310/webmake/vars"
	"io"
	"encoding/json"
)

type Scripts []map[string]map[string]string;

func GetScripts(ctx context.Context) Scripts {
	if req, err := http.NewRequestWithContext(
		ctx, 
		"GET", 
		"https://raw.githubusercontent.com/DOMIN1310/webmake/master/res/scripts.json",
		nil,
	); err != nil {
		log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, "unable to get scripts request");
		return Scripts{};		
	} else {
		if res, err := http.DefaultClient.Do(req); err != nil {
			log.Fatalf("%v:%v%v\n", v.ERROR, v.RESET, "unable to get response from the request to get scripts");
			return Scripts{};
		} else {
			if buffer, err := io.ReadAll(res.Body); err != nil {
				res.Body.Close();
				log.Printf("%v:%v%v\n", v.ERROR, v.RESET, "unable to read response");
				return Scripts{};
			} else {
				defer res.Body.Close();
				var variable Scripts;
				if err := json.Unmarshal(buffer, &variable); err != nil {
					log.Printf("%v:%v%v\n", v.ERROR, v.RESET, "unable to marshal data");
					return Scripts{};
				} else {
					return variable;
				}
			}
		}
	}
}

func GetURLBody(url string, ctx context.Context, method string, file string) map[string]string {
	if req, err := http.NewRequestWithContext(ctx, method, url, nil); err != nil {
		log.Printf("%v:%v%v\n", v.WARN, v.RESET, "unable to get any request within context and destiny url");
		return map[string]string{file: ""};	
	} else if ctx.Err() != nil {
		log.Printf("%v:%v%v\n", v.WARN, v.RESET, "unable to complete request: context error")
		return map[string]string{file: ""};
	} else {
		if res, err := http.DefaultClient.Do(req); err != nil {
			res.Body.Close();
			log.Printf("%v:%v%v\n", v.WARN, v.RESET, "unable to complete request with proper response");
			return map[string]string{file: ""};
		} else {
			defer res.Body.Close();
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