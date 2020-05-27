package foo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

var fooMap = struct {
	sync.RWMutex
	m map[int]Foo
}{m: make(map[int]Foo)}

func init() {
	fmt.Println("loading foos...")
	prodMap, err := loadFooMap()
	fooMap.m = prodMap

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d foos loaded...\n", len(fooMap.m))
}

func loadFooMap() (map[int]Foo, error) {
	fileName := "foos.json"
	_, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}

	file, _ := ioutil.ReadFile(fileName)
	fooList := make([]Foo, 0)
	err = json.Unmarshal([]byte(file), &fooList)

	if err != nil {
		log.Fatal(err)
	}

	prodMap := make(map[int]Foo)
	for i := 0; i < len(fooList); i++ {
		prodMap[fooList[i].ProductID] = fooList[i]
	}

	return prodMap, nil
}

func getFoo(productID int) *Foo {
	fooMap.RLock()
	defer fooMap.RUnlock()

	if foo, ok := fooMap.m[productID]; ok {
		return &foo
	}

	return nil
}

func removeFoo(productID int) {
	fooMap.RLock()
	defer fooMap.RUnlock()

	delete(fooMap.m, productID)
}

func getFooList() []Foo {
	fooMap.RLock()
	foos := make([]Foo, 0, len(fooMap.m))

	for _, value := range fooMap.m {
		foos = append(foos, value)
	}

	fooMap.RUnlock()

	return foos
}

func getFooIds() []int {
	fooMap.RLock()
	fooIds := []int{}

	for key := range fooMap.m {
		fooIds = append(fooIds, key)
	}

	fooMap.RUnlock()
	sort.Ints(fooIds)

	return fooIds
}

func getNextFooID() int {
	fooIDs := getFooIds()
	return fooIDs[len(fooIDs)-1] + 1
}

func addOrUpdateFoo(foo Foo) (int, error) {
	addOrUpdateID := -1

	if foo.ProductID > 0 {
		oldFoo := getFoo(foo.ProductID)

		if oldFoo == nil {
			return 0, fmt.Errorf("Foo id [%d] doesn't exist", foo.ProductID)
		}

		addOrUpdateID = foo.ProductID
	} else {
		addOrUpdateID = getNextFooID()
		foo.ProductID = addOrUpdateID
	}

	fooMap.Lock()
	fooMap.m[addOrUpdateID] = foo
	fooMap.Unlock()

	return addOrUpdateID, nil
}
