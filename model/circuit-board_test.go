package model

import (
	"reflect"
	"testing"
)

func TestCircuitBoardFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *CircuitBoard
		wantErr bool
	}{
		// Happy path tests
		{
			name: "Simplest",
			args: args{
				s: "|ST|",
			},
			want: &CircuitBoard{
				width:    1,
				height:   1,
				chips:    []Chip{Chip(StartChip)},
				startPos: Position{X: 0, Y: 0},
			},
		},
		{
			name: "One row Forward",
			args: args{
				s: "|ST -> MF|",
			},
			want: &CircuitBoard{
				width:    2,
				height:   1,
				chips:    []Chip{Chip(StartChip).WithArrowYes(East), Chip(ForwardChip)},
				startPos: Position{X: 0, Y: 0},
			},
		},
		{
			name: "Two rows Turn Left",
			args: args{
				s: `
|TL|
| ^|
|ST|`,
			},
			want: &CircuitBoard{
				width:    1,
				height:   2,
				chips:    []Chip{Chip(TurnLeftChip), Chip(StartChip).WithArrowYes(North)},
				startPos: Position{X: 0, Y: 1},
			},
		},
		{
			name: "3x2 with yes and no",
			args: args{
				s: `
|ST -> W? y> TL|
|      nv     v|
|      MF <- ..|`,
			},
			want: &CircuitBoard{
				width:  3,
				height: 2,
				chips: []Chip{
					// Row 1
					Chip(StartChip).WithArrowYes(East),
					Chip(IsWallAheadChip).WithArrowYes(East).WithArrowNo(South),
					Chip(TurnLeftChip).WithArrowYes(South),

					// Row 2
					Chip(NoChip),
					Chip(ForwardChip),
					Chip(NoChip).WithArrowYes(West),
				},
				startPos: Position{X: 0, Y: 0},
			},
		},
		// TODO: Add sad path test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CircuitBoardFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("CircuitBoardFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CircuitBoardFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
