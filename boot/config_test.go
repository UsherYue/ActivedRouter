package boot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

//load config test
func TestConfig(t *testing.T) {
	file, _ := os.Open("../config/client.json")
	bts, _ := ioutil.ReadAll(file)
	var server ClientConfigData
	json.Unmarshal(bts, &server)
	fmt.Println(server)
	//	var x map[string]string = map[string]string{"x": "x"}
	//	fmt.Println((x))
}
