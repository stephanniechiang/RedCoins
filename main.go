package main

func main() {
    a := App{} 
    // You need to set your Username and Password here
    a.Initialize("root", "moedasvermelhas", "RedCoins")

    a.Run(":8080")
}