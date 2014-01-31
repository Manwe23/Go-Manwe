// TheGameMap project TheGameMap.go
package TheGameMap

import (
	"database/sql"
	"fmt"
	_ "github.com/ziutek/mymysql/godrv"
	"log"
	"math/rand"
	"time"
)

type Mapa struct {
}

func (m Mapa) Render() {
	fmt.Println("Rendering...")
}

func (m Mapa) Load() {
	fmt.Println("preparing...")
	con, err := sql.Open("mymysql", "sql.pmp31121984.nazwa.pl:3307/pmp31121984_22/pmp31121984_22/StudenT23@")
	fmt.Println("passed")

	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()
}

func (m Mapa) Get(long int, lat int) {
	long, lat = m.GetRandom()
	fmt.Printf("[%d,%d]=H2C1T5S7\n", long, lat)
}

func (m Mapa) Set(long int, lat int, attr_name string, attr_value string) {
	fmt.Printf("[%d,%d].%s = %s\n", long, lat, attr_name, attr_value)
}

func (m Mapa) GetRandom() (long int, lat int) {
	//var t time.Time
	rand.Seed(time.Now().Unix())
	long = 180 - rand.Intn(360)
	lat = 180 - rand.Intn(360)
	return
}
