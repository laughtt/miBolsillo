package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type geometry interface {
	area(w int , h int)int
	shape(s string)string
}

type circle struct{
	r int
}
func (c *circle) area(a int , b int) int {
	return 1
}
func (c *circle) shape(string)string {
	return "a"
}

func measure(g geometry) int {
	return g.area(1 , 2)
}
func main(){
measure(circle{r : 1})
}
