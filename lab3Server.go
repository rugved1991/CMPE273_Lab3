package main
import  ("github.com/julienschmidt/httprouter"
    "fmt"
    "net/http"
    "strconv"
    "encoding/json"
    "strings"
    "sort")


type KeyVals struct{
  Key int `json:"key,omitempty"`
  Val string  `json:"value,omitempty"`
} 


var server1,server2,server3 [] KeyVals
var index1,index2,index3 int
type ByKey []KeyVals
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

func main(){
  index1 = 0
  index2 = 0
  index3 = 0
  mux := httprouter.New()
    mux.GET("/keys",getAllKeys)
    mux.GET("/keys/:key_id",GetKey)
    mux.PUT("/keys/:key_id/:value",PutKeys)
    go http.ListenAndServe(":3000",mux)
    go http.ListenAndServe(":3001",mux)
    go http.ListenAndServe(":3002",mux)
    select {}
}

func GetKey(respWrit http.ResponseWriter, request *http.Request,param httprouter.Params){ 
  out := server1
  ind := index1
  port := strings.Split(request.Host,":")
  if(port[1]=="3001"){
    out = server2 
    ind = index2
  }else if(port[1]=="3002"){
    out = server3
    ind = index3
  } 
  key,_ := strconv.Atoi(param.ByName("key_id"))
  for i:=0 ; i< ind ;i++{
    if(out[i].Key==key){
      res,_:= json.Marshal(out[i])
      fmt.Fprintln(respWrit,string(res))
    }
  }
}

func PutKeys(respWrit http.ResponseWriter, request *http.Request,param httprouter.Params){
  port := strings.Split(request.Host,":")
  key,_ := strconv.Atoi(param.ByName("key_id"))
  if(port[1]=="3000"){
    server1 = append(server1,KeyVals{key,param.ByName("value")})
    index1++
  }else if(port[1]=="3001"){
    server2 = append(server2,KeyVals{key,param.ByName("value")})
    index2++
  }else{
    server3 = append(server3,KeyVals{key,param.ByName("value")})
    index3++
  } 
}

func getAllKeys(respWrit http.ResponseWriter, request *http.Request,param httprouter.Params){
  port := strings.Split(request.Host,":")
  if(port[1]=="3000"){
    sort.Sort(ByKey(server1))
    res,_:= json.Marshal(server1)
    fmt.Fprintln(respWrit,string(res))
  }else if(port[1]=="3001"){
    sort.Sort(ByKey(server2))
    res,_:= json.Marshal(server2)
    fmt.Fprintln(respWrit,string(res))
  }else{
    sort.Sort(ByKey(server3))
    res,_:= json.Marshal(server3)
    fmt.Fprintln(respWrit,string(res))
  }
}







