package main

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

// Função que gera todas as combinações possíveis de uma string de comprimento n usando os caracteres fornecidos.
func generateText(chars string, length int) []string {
    var combinations []string
    var generate func(prefix string, length int)
    
    generate = func(prefix string, length int) {
        if length == 0 {
            combinations = append(combinations, prefix)
            return
        }
        for _, c := range chars {
            generate(prefix+string(c), length-1)
        }
    }
    
    generate("", length)
    return combinations
}

func textToMD5(text string) string {
    hash := md5.Sum([]byte(text))
    return hex.EncodeToString(hash[:])
}

// Função que tenta encontrar a senha correta usando força bruta.
func singleProcess(initialText, chars string, length int, pwd string, flag *int32, wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 1; i <= length; i++ {
        for _, text := range generateText(chars, i) {
            combinedText := initialText + text
            if atomic.LoadInt32(flag) == 1 {
                return
            }
            if textToMD5(combinedText) == pwd {
                atomic.StoreInt32(flag, 1)
                fmt.Printf("A senha é %s\n", combinedText)
                return
            }
        }
    }
}

func main() {
    var wg sync.WaitGroup
    var flag int32 = 0
    chars := "abcdefghijklmnopqrstuvwxyz"
    //length := 5
    var length int
    fmt.Print("Digite o tamanho da senha: ")
    fmt.Scanln(&length)
    //pwd := "5d41402abc4b2a76b9719d911017c592" // Exemplo de hash MD5 para "hello"
    // Solicita o hash MD5 ao usuário
    var pwd string
    fmt.Print("Digite o hash MD5: ")
    fmt.Scanln(&pwd)

    // Inicia a medição do tempo
    start := time.Now()

    for _, c := range chars {
        wg.Add(1)
        go singleProcess(string(c), chars, length, pwd, &flag, &wg)
    }

    wg.Wait()

    // Calcula o tempo de processamento
    elapsed := time.Since(start)

    if flag == 1 {
        fmt.Println("Senha quebrada com sucesso!")
    } else {
        fmt.Println("Não foi possível quebrar a senha.")
    }

    // Exibe o tempo de processamento
    fmt.Printf("Tempo de processamento: %s\n", elapsed)
}