package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Struct para ejemplo
type Person struct {
	Name string
	Age  int
}

// 1. Devuelve error
func returnError() error {
	return errors.New("algo salió mal")
}

// 2. Devuelve bool
func returnBool() bool {
	return true
}

// 3. Devuelve int
func returnInt() int {
	return 42
}

// 4. Devuelve string
func returnString() string {
	return "Hola mundo"
}

// 5. Devuelve struct
func returnStruct() Person {
	return Person{Name: "Juan", Age: 30}
}

// 6. Devuelve array de int
func returnIntArray() []int {
	return []int{1, 2, 3, 4, 5}
}

// 7. Devuelve array de string
func returnStringArray() []string {
	return []string{"uno", "dos", "tres"}
}

// 8. Devuelve JSON string
func returnJSON() string {
	p := Person{Name: "Ana", Age: 25}
	data, _ := json.Marshal(p)
	return string(data)
}

// 9. Entrada string y devuelve int
func lengthOfString(s string) int {
	return len(s)
}

// 10. Entrada int y devuelve bool
func isEven(n int) bool {
	return n%2 == 0
}

// 11. Entrada struct y devuelve string
func greetPerson(p Person) string {
	return fmt.Sprintf("Hola, %s. Tienes %d años.", p.Name, p.Age)
}

// 12. Entrada array y devuelve suma
func sumArray(arr []int) int {
	sum := 0
	for _, n := range arr {
		sum += n
	}
	return sum
}

// 13. Entrada string y devuelve error o nil
func validateString(s string) error {
	if s == "" {
		return errors.New("cadena vacía")
	}
	return nil
}

// 14. Entrada y salida múltiples
func complexFunction(s string, n int) (string, int, error) {
	if s == "" {
		return "", 0, errors.New("string vacío")
	}
	return s, n * 2, nil
}

// 15. Entrada bool y devuelve struct
func flagToPerson(flag bool) Person {
	if flag {
		return Person{Name: "Verdadero", Age: 1}
	}
	return Person{Name: "Falso", Age: 0}
}

// 16. Entrada y salida JSON
func modifyJSON(jsonStr string) (string, error) {
	var p Person
	err := json.Unmarshal([]byte(jsonStr), &p)
	if err != nil {
		return "", err
	}
	p.Age += 1
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// main
/*
func main() {
	fmt.Println(returnError())
	fmt.Println(returnBool())
	fmt.Println(returnInt())
	fmt.Println(returnString())
	fmt.Println(returnStruct())
	fmt.Println(returnIntArray())
	fmt.Println(returnStringArray())
	fmt.Println(returnJSON())
	fmt.Println(lengthOfString("Hola"))
	fmt.Println(isEven(4))
	fmt.Println(greetPerson(Person{"Luis", 40}))
	fmt.Println(sumArray([]int{1, 2, 3}))
	fmt.Println(validateString(""))
	fmt.Println(complexFunction("Go", 5))
	fmt.Println(flagToPerson(true))

	jsonOut, err := modifyJSON(`{"Name":"Pepe","Age":20}`)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(jsonOut)
	}
}

*/
