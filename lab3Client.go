package main  

  
import (  
    "fmt"  
    "hash/crc32"  
    "sort"     
    "net/http"
    "encoding/json" 
    "io/ioutil"
    "os"
    "strings"
)  
   
type HashCircle []uint32  

type KeyVals struct{
    Key int `json:"key,omitempty"`
    Val string `json:"value,omitempty"`
}



func (hr HashCircle) Len() int {  
    return len(hr)  
}  
  
func (hr HashCircle) Less(i, j int) bool {  
    return hr[i] < hr[j]  
}  
  
func (hr HashCircle) Swap(i, j int) {  
    hr[i], hr[j] = hr[j], hr[i]  
}  
  
type Node struct {  
    Id       int  
    IP       string    
}  
  

  
type ConstHash struct {  
    Nodes       map[uint32]Node  
    IsPresent   map[int]bool  
    Circle      HashCircle  
    
}  

func main() {   
    circle := NewConsistentHash()      
    circle.AddNode(NewNode(0, "127.0.0.1:3000"))
    circle.AddNode(NewNode(1, "127.0.0.1:3001"))
    circle.AddNode(NewNode(2, "127.0.0.1:3002")) 
    
    if(os.Args[1]=="PUT"){
        key := strings.Split(os.Args[2],"/")
        PutKey(circle,key[0],key[1])
    } else if ((os.Args[1]=="GET") && len(os.Args)==3){
        GetKey(os.Args[2],circle)
    } else {
        GetAll("http://127.0.0.1:3000/keys")
        GetAll("http://127.0.0.1:3001/keys")
        GetAll("http://127.0.0.1:3002/keys")
    }
}  

func (hr *ConstHash) Get(key string) Node {  
    hash := hr.GetHashValue(key)  
    i := hr.SearchForNode(hash)  
    return hr.Nodes[hr.Circle[i]]  
}  
  
func (hr *ConstHash) AddNode(node *Node) bool {  
 
    if _, ok := hr.IsPresent[node.Id]; ok {  
        return false  
    }  
    str := hr.ReturnNodeIP(node)  
    hr.Nodes[hr.GetHashValue(str)] = *(node)
    hr.IsPresent[node.Id] = true  
    hr.SortHashCircle()  
    return true  
}    

func NewNode(id int, ip string) *Node {  
   return &Node{  
       Id:       id,  
       IP:       ip,  
   }  
}  

func NewConsistentHash() *ConstHash {  
    return &ConstHash{  
        Nodes:     make(map[uint32]Node),   
        IsPresent: make(map[int]bool),  
        Circle:      HashCircle{},  
    }  
}  
  
func (hr *ConstHash) SearchForNode(hash uint32) int {  
    i := sort.Search(len(hr.Circle), func(i int) bool {return hr.Circle[i] >= hash })  
    if i < len(hr.Circle) {  
        if i == len(hr.Circle)-1 {  
            return 0  
        } else {  
            return i  
        }  
    } else {  
        return len(hr.Circle) - 1  
    }  
} 
  
func (hr *ConstHash) SortHashCircle() {  
    hr.Circle = HashCircle{}  
    for k := range hr.Nodes {  
        hr.Circle = append(hr.Circle, k)  
    }  
    sort.Sort(hr.Circle)  
}  
 
func GetKey(key string,circle *ConstHash){
    var out KeyVals 
    ipAddress:= circle.Get(key)
    address := "http://"+ipAddress.IP+"/keys/"+key
    fmt.Println(address)
    response,err:= http.Get(address)
    if err!=nil{
        fmt.Println("Error:",err)
    }else{
        defer response.Body.Close()
        contents,err:= ioutil.ReadAll(response.Body)
        if(err!=nil){
            fmt.Println(err)
        }
        json.Unmarshal(contents,&out)
        result,_:= json.Marshal(out)
        fmt.Println(string(result))
    }
}

func (hr *ConstHash) ReturnNodeIP(node *Node) string {  
    return node.IP 
}  
  
func (hr *ConstHash) GetHashValue(key string) uint32 {  
    return crc32.ChecksumIEEE([]byte(key))  
}  
  


 
  
func PutKey(circle *ConstHash, str string, input string){
        ipAddress := circle.Get(str)  
        address := "http://"+ipAddress.IP+"/keys/"+str+"/"+input
		fmt.Println(address)
        req,err := http.NewRequest("PUT",address,nil)
        client := &http.Client{}
        resp, err := client.Do(req)
        if err!=nil{
            fmt.Println("Error:",err)
        }else{
            defer resp.Body.Close()
            fmt.Println("PUT Request successfully completed")
        }  
}  



func GetAll(address string){
     
    var out []KeyVals
    response,err:= http.Get(address)
    if err!=nil{
        fmt.Println("Error:",err)
    }else{
        defer response.Body.Close()
        contents,err:= ioutil.ReadAll(response.Body)
        if(err!=nil){
            fmt.Println(err)
        }
        json.Unmarshal(contents,&out)
        result,_:= json.Marshal(out)
        fmt.Println(string(result))
    }
}
