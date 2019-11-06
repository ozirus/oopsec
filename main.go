package main

import (
	"log"
	"github.com/jpillora/go-tld"
	"strings"
	"encoding/csv"
	"os"
	"bufio"
	"fmt"
	"regexp"
	"github.com/likexian/whois-go"
	"github.com/likexian/whois-parser-go"
)

var iocler []string

func main() {

	file, err := os.Open("iocs")
	if err != nil {
	    log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
        // Line'lar her zaman strings.ToLower geçmeli parametre olurken	
	for scanner.Scan() {
	    fmt.Print(temizle(scanner.Text())+"\t\t\t\t\t\t\t\t|||-->")
	    fmt.Println(checkIOCType(temizle(scanner.Text())))
	    //domainTop1Mmi(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
	    log.Fatal(err)
	}
}


func temizle(_domain string) string{
	_domainL := strings.ToLower(_domain)
	bir:= strings.Replace(_domainL, "[.]",".",-1)
	iki:= strings.Replace(bir,"hxxp","http",-1)
	son := extractDomainFromURL(iki)

	return son
}

func domainTop1Mmi(_domain string) {
	result, err := whois.Whois(_domain)
	if err == nil {
		resultik, erriki := whoisparser.Parse(result)
		if erriki == nil {
			createdate := resultik.Registrar.CreatedDate
			fmt.Print(string(createdate))

			_temizDomain := strings.TrimRight(_domain,"\n")
			 lines, err := readCsv("top-1m.csv")
			 if err != nil {
					panic(err)
		         }

			for _, line := range lines {
				data := csvLine{
				    sira: line[0],
				    domain: line[1],
				}
				if _temizDomain==data.domain {
					fmt.Print(" -> [!] Top 1M'da")
				}
	                 }
			 fmt.Println()


			}
	}
}
func check(e error) {
    if e != nil {
        panic(e)
    }
}

type csvLine struct {
    sira string
    domain string
}

func readCsv(filename string) ([][]string, error) {

    f, err := os.Open(filename)
    if err != nil {
        return [][]string{}, err
    }
    defer f.Close()

    lines, err := csv.NewReader(f).ReadAll()
    if err != nil {
        return [][]string{}, err
    }

    return lines, nil
}

func checkIOCType(line string) string {
	reIP := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
        reDomain := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z
 ]{2,3})$`)

	if reDomain.MatchString(line) {
		return "domain"
	} else if reIP.MatchString(line) {
		return "ip"
	}
		return "hash"
}

func extractDomainFromURL(gelenURL string) string {
if strings.Contains(gelenURL,"/") {
	    u,_ := tld.Parse(gelenURL)	
	if u.Subdomain == "" {
		return u.Domain+"."+u.TLD
        }	
	return u.Subdomain+"."+u.Domain+"."+u.TLD
}
return gelenURL

}
