package main

import(
	"fmt"
	"os"
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"container/list"
	"errors"
	"regexp"
)

func precompute_jumps(input []byte) (map[uint]uint, error){
	stack := list.New()
	var (
		jumps map[uint]uint = make(map[uint]uint)

		plen uint = uint(len(input))
		fpos uint = 0
	)

	for fpos < plen{
		switch input[fpos]{
		case '[':
			stack.PushBack(fpos)
		case ']':
			if stack.Len() == 0{
				return nil, errors.New("Unexpected error with stack length")
			}
			tget := stack.Remove(stack.Back()).(uint)

			jumps[tget] = fpos
			jumps[fpos] = tget
		}
		fpos++
	}
	if stack.Len() != 0{
		return nil, errors.New("Too many opening brackets")
	}

	return jumps, nil
}

func interpret(input []byte, j io.Reader, w io.Writer){
	reader := bufio.NewReader(j)
	list := make([]byte, 30000)
	cnt := 0	

	jumps, err := precompute_jumps(input)
	if err != nil{
		log.Fatal(err)
	}
	
	var i uint = 0
	var progLen uint = uint(len(input))
	stringInput := string(input)

	for i < progLen{
		switch chara := stringInput[i]; chara{
		case '+':
			list[cnt]++
		case '-':
			list[cnt]--	
		case '<':
			cnt--
			if cnt < 0{
				cnt = 29999
			}
		case '>':
			cnt++
			if cnt > 29999{
				cnt = 0
			}
		case '.':
			fmt.Fprintf(w, "%c", list[cnt])
		case ',':
			if list[cnt], err = reader.ReadByte(); err != nil{
				os.Exit(0)
			}
		case '[':
			if list[cnt] == 0{
				i = jumps[i]
			}
		case ']':
			if list[cnt] != 0 {
				i = jumps[i]
			}
		}
		i++
	}
}


func main(){
	osArgs := os.Args[1]
	r, err := ioutil.ReadFile(osArgs)
	if err != nil{
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	re := regexp.MustCompile(`\r?\n| |[a-zA-Z0-9]`)
	bytes := re.ReplaceAll(r, []byte{})
	
	interpret(bytes, os.Stdin, os.Stdout)
}
