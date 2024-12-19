package main

import (
    "github.com/alirezadp10/hokm/cmd"
    "github.com/joho/godotenv"
)

func main() {
    _ = godotenv.Load()
    cmd.Execute()

    //sqliteClient := sqlite.GetNewConnection()
    //
    //fmt.Printf("=======\n")
    //if gid, ok := sqlite.DoesPlayerHaveAnActiveGame(sqliteClient, "amir"); ok {
    //    fmt.Printf("..........\n")
    //    fmt.Printf("%v\n", *gid)
    //}
}
