# Glicko2 implementation for Golang

This is an implementation of the [Glicko2](http://www.glicko.net/glicko/glicko2.pdf) rating system for the [Go](https://golang.org/) programming language.

Example usage:

```  
  player1 := &Rating{1500.0, 200.0, 0.06, []Result{}}
  player2 := &Rating{1400.0, 30.0, 0.06, []Result{}}
  player3 := &Rating{1550.0, 100.0, 0.06, []Result{}}
  player4 := &Rating{1700.0, 300.0, 0.06, []Result{}}

  player1.AddResult(Result{player2, 1.0})
  player2.AddResult(Result{player1, 0.0})

  player1.AddResult(Result{player3, 0.0})
  player3.AddResult(Result{player1, 1.0})

  player1.AddResult(Result{player4, 0.0})
  player4.AddResult(Result{player1, 1.0})

  s := &System{Tau: 0.5, Mu: 1500, Phi: 350, Epsilon: 0.000001, Players: []*Rating{player1, player2, player3, player4}}

  s.Update()
```
