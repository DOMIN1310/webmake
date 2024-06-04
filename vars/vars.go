package vars

type Template struct{
	Dir string `json:"dir"`;
	Findex string `json:"findex"`;
	Tmplindex string `json:"tmplindex"`;
	Styleindex string `json:"styleindex"`;
	Git bool `json:"git"`;
}

var (
	ERROR string = "\033[31m[ERROR]";
	WARN string = "\033[33m[WARN]";
	SUCCESS string = "\033[32m[SUCCESS]"
	RESET string = "\033[0m"
)