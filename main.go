package main

func main() {
    a := App{} 
    // You need to set your Username and Password here
    a.Initialize("redcoins", "moedasvermelhas", "rest_api_example")

    a.Run(":8080")
}