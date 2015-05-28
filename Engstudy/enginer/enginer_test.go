package enginer

import (
	"fmt"
	//"os"
	"testing"
	"time"
)

func hs(job *Job) error {
	fmt.Printf("%#v\n", job)
	return nil
}

func TestEnginer(t *testing.T) {
	eng := New()
	eng.Register("hs", hs)
	eng.Onshutdown(func() {
		fmt.Println("HEHE4")
		time.Sleep(time.Second * 2)
	})
	j := eng.Job("hs", "cd")
	j.SetEnv("he", "meili")
	fmt.Println(j.GetEnvString("he"))
	//j.Stdout.Add(os.Stdout)
	fmt.Println("HEHE3")
	j.Run()
	eng.Shutdown()
	fmt.Println("HEHE5")
}
