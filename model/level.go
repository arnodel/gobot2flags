package model

type Level struct {
	Name        string
	Maze        *Maze
	BoardWidth  int
	BoardHeigth int
}

func LevelFromString(name string, s string) (*Level, error) {
	m, err := MazeFromString(s)
	if err != nil {
		return nil, err
	}
	return &Level{
		Name:        name,
		Maze:        m,
		BoardWidth:  9,
		BoardHeigth: 9,
	}, nil
}

func (l *Level) BoardSize() (int, int) {
	return l.BoardWidth, l.BoardHeigth
}
