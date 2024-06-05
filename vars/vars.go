package vars

type Template struct{
	Findex string `json:"findex"`;
	Tmplindex string `json:"tmplindex"`;
	Styleindex string `json:"styleindex"`;
	Git bool `json:"git"`;
}

var (
	FINISHED string = "\033[1;37m[FINISHED]";
	ERROR string = "\033[1;31m[ERROR]";
	WARN string = "\033[1;33m[WARN]";
	SUCCESS string = "\033[1;32m[SUCCESS]";
	RESET string = "\033[0m";
	INIT string = "\033[1;36m[INIT]";
	DONE string = "\033[1;95m[DONE]";
	CREATION string = "\033[1;34m[CREATED]"
)