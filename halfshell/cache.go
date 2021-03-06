package halfshell

import (
	"os"
	"time"
	"fmt"
	"io/ioutil"
	"strings"
	"os/exec"
)

const TMP_FOLDER string = "/tmp/halfshell"

func CacheInit() {

	if _, err := os.Stat(TMP_FOLDER); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(TMP_FOLDER, 0777); err != nil {
				panic(err)
			}

		} else {
			panic(err)
		}
	}
}


func CacheRead(path string) (*Image, error) {

	start := time.Now()
	content, err := os.Open(TMP_FOLDER+"/"+path)
	defer content.Close()
	if err == nil {
		image, err := NewImageFromBuffer(content)
		if err != nil{
			fmt.Printf("Problem creating image from local file %v\n", path)
			return nil, err
		} else {
			fmt.Printf("Successfully retrieved image from local cache: %v %v\n", path, time.Since(start))
			return image, nil
		}
	} else {
		return nil, err
	}

}

func CacheWrite(path string, img *Image)  {

	start := time.Now()
	image_bytes,num := img.GetBytes()
	if num > 0 {
		err := ioutil.WriteFile(TMP_FOLDER+"/"+path, image_bytes, 0644)
		if err != nil {
			fmt.Printf("Problem on updating cache %v\n", path)
		}
	}

	fmt.Printf("Successfully update cache: %v %v\n", path, time.Since(start))
}

func CleanCacheFolder()  {

	cmd := `find `+TMP_FOLDER+` -type d -ctime +3 -exec rm -rf {} \;`
	fmt.Println(cmd)

	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head,parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)

}

func CacheDeamonRun() {

	for {
		CleanCacheFolder()
		time.Sleep(1*time.Hour)
	}
}
