package main

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "os"
    "sync"
    "bufio"
)

func textToMD5(text string) string {
    hash := md5.Sum([]byte(text))
    return hex.EncodeToString(hash[:])
}

func generateText(n int, chars string) []string {
    if n == 0 {
        return []string{""}
    }
    pwList := generateText(n-1, chars)
    var result []string
    for _, pw := range pwList {
        for _, c := range chars {
            result = append(result, pw+string(c))
        }
    }
    return result
}

func singleProcess(initialText, chars string, length int, pwd string, flag *int32, wg *sync.WaitGroup, verbose bool) {
    defer wg.Done()
    for i := 1; i <= length; i++ {
        for _, text := range generateText(i, chars) {
            combinedText := initialText + text
            if verbose {
                fmt.Printf("%d is trying %s\n", os.Getpid(), combinedText)
            }
            if *flag == 1 {
                return
            }
            if textToMD5(combinedText) == pwd {
                *flag = 1
                fmt.Printf("The password is %s\n", combinedText)
                return
            }
        }
    }
}

func main() { // Leitura do hash MD5 da senha
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Digite o hash MD5 da senha: ")
    pwd, _ := reader.ReadString('\n')
    pwd = pwd[:len(pwd)-1] // Remove o caractere de nova linha

    // Leitura do comprimento da senha
    fmt.Print("Digite o comprimento da senha: ")
    var length int
    fmt.Scanf("%d", &length)

    chars := "abcdefghijlmnopqrstuvwxzABCDEFGHIJLMNOPQRSTUVWXZ0123456789#$%&*+-.*="
    var wg sync.WaitGroup
    var flag int32 = 0
    verbose := true

    for _, c := range chars {
        wg.Add(1)
        go singleProcess(string(c), chars, length, pwd, &flag, &wg, verbose)
    }

    wg.Wait()

    if flag == 1 {
        fmt.Println("Senha quebrada com sucesso!")
    } else {
        fmt.Println("Não foi possível quebrar a senha.")
    }
}