package service

import (
    "fmt"
    "github.com/alirezadp10/hokm/pkg/mocks"
    "testing"
)

func TestChooseFirstKing(t *testing.T) {
    ps := NewPlayersService(new(mocks.PlayersRepositoryContract))
    x, y := ps.ChooseFirstKing()
    fmt.Println(x)
    fmt.Println(y)
}
