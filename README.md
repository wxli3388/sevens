# Raw of sevens poker game with websocket.
The game is designed to play row of sevens. 
By default, the game involves three computer players upon joining the game room. However, it can be modified to allow other users to join and play.

# To-Do List
1. Calculate and summarize the scores at the end of each game.
2. Enhance the user interface (UI) for better user experience.
3. Resolve the issue related to writing to a closed channel when a user disconnects from the WebSocket.
4. Optimize the card comparison logic for improved gameplay.


# How to use
Start go server
```
go run cmd/main.go 
```

Open the html 
```
sevens/web/ws
```


# demo

![image](https://github.com/wxli3388/sevens/blob/main/demo.gif)
