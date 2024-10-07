# Remove the Islands Case with Go

It's been some time since I first started studying Go. Coming from a long time using Javascript for automating some tasks and creating services at my current job, my first contact with Go was simply delightful.

No types cohersion, no overnight revolution with a brand new framework, a formidable standard library... Everything looks so simple and yet so powerful.

So, in order to share my jorney learning Go and getting better with data structures and algorithms, I'll start posting somethings I've been trying with Go.

### The "Remove the Islands" case

Some time ago I came accross this interview case calsse Remove the Islands: given a grid of 1 and 0 in which 1 represents land and 0 represents water, you should remove from the grid all islands, pieces of land not connected to the border. It's important to notice that elements cannot connect diagonally.

This is the input for the case:

```javascript
[
    [1, 0, 0, 0, 0, 0],
    [0, 1, 0, 1, 1, 1],
    [0, 0, 1, 0, 1, 0],
    [1, 1, 0, 0, 1, 0],
    [1, 0, 1, 1, 0, 0],
    [1, 0, 0, 0, 0, 1]
]
```

And this is the expected output:

```javascript
[
    [1, 0, 0, 0, 0, 0],
    [0, 0, 0, 1, 1, 1],
    [0, 0, 0, 0, 1, 0],
    [1, 1, 0, 0, 1, 0],
    [1, 0, 0, 0, 0, 0],
    [1, 0, 0, 0, 0, 1]
]
```

### My solution for the case

Using Go, I've developed a solution exploring the concepts of recursion and memoization from dynamic programming.

To setup memoization, I've created two custom types: _Memo_ using a hash map where the keys represent the element's coordinates and _Board_ which represents the grid and initialized the memo variable.

```go
type (
    Board [][]int
    Memo  map[[2]int]bool
)

var memo = Memo{}
```

After that, I've declared a function responsible for checking if a given element, represented by it's coordinates _row_ and _col_, is connected to the border of the _Board_. The _increment_ parameter is used to offset the element, checking for its surroundings up, down, left and right. As neighbor elements might check the same surroundings, it's always important to check if the element to be verified is already in _memo_.

```go
func checkDirection(board Board, row, col int, increment map[byte]int) bool {
    if _, ok := memo[[2]int{row, col}]; ok {
    } else if col == 0 || row == 0 || col == len(board[row])-1 || row == len(board)-1 {
        memo[[2]int{row, col}] = board[row][col] == 1
    } else if board[row][col] == 0 {
        memo[[2]int{row, col}] = false
    } else {
        col, row = col+increment['x'], row+increment['y']
        return checkDirection(board, row, col, increment)
    }
    return memo[[2]int{row, col}]
}
```

The _checkDirection_ function allows to traverse the board, checking for each element if it's possible to find a path connecting the given element to the board passing only throgh land elements, which are represented by 1. The _connectsToBoard_ implements _checkDirection_ for each possible path to the border: up, down, left and right.

```go
func connectsToBorder(board Board, row, col int) bool {
    if v, ok := memo[[2]int{row, col}]; ok {
        return v
    }
    var connects bool
    directions := []map[byte]int{
        {'x': -1, 'y': 0},
        {'x': 1, 'y': 0},
        {'x': 0, 'y': -1},
        {'x': 0, 'y': 1},
    }
    for _, increment := range directions {
        connects = connects || checkDirection(board, row, col, increment)
    }
    memo[[2]int{row, col}] = connects
    return connects
}
```

If _connectsToBorder_ returns true, it means there is at least one path of land connecting the element at position col, row to the border of the board.

Now it's possible to remove all islands, traversing the board and replacing every island with water, represented by 0. The function _removeIslands_ is responsible for removing every island from the board:

```go
func removeIslands(board Board) Board {
    result := make(Board, len(board))
    for row := 0; row < len(board); row++ {
        result[row] = make([]int, len(board[row]))
        for col := 0; col < len(board[row]); col++ {
            connects := connectsToBorder(board, row, col)
            if board[row][col] == 1 && !connects {
                result[row][col] = 0
            } else {
                result[row][col] = board[row][col]
            }
        }
    }
    return result
}
```

Using memoization was crucial for this case, in order to achieve a time complexity of O(n*m), where _n_ and _m_ represents the dimensions of the board.

If you are interested, you can check out [my solution on github](https://github.com/gzitei/just_studying/blob/main/remove-the-islands/). I'm really interested in hearing suggestions and other aproaches to solve the problem, please send me a DM with your ideas!

> [!NOTE]
> TESTE!!!
