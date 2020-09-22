package main

import (
	"testing"
)

func TestChip_Type(t *testing.T) {
	tests := []struct {
		name string
		c    Chip
		want ChipType
	}{
		{
			name: "WithType",
			c:    Chip(0).WithType(TurnLeftChip),
			want: TurnLeftChip,
		},
		{
			name: "WithType twice",
			c:    Chip(0).WithType(ForwardChip).WithType(IsFloorBlueChip),
			want: IsFloorBlueChip,
		},
		{
			name: "WithArrow",
			c:    Chip(0).WithType(ForwardChip).WithArrowYes(North),
			want: ForwardChip,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Type(); got != tt.want {
				t.Errorf("Chip.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChip_ArrowYes(t *testing.T) {
	tests := []struct {
		name  string
		c     Chip
		want  Orientation
		want1 bool
	}{
		{
			name:  "Default",
			c:     Chip(0),
			want1: false,
		},
		{
			name:  "WithArrowYes",
			c:     Chip(0).WithArrowYes(East),
			want:  East,
			want1: true,
		},
		{
			name:  "WithArrowYes ClearArrowYes",
			c:     Chip(0).WithArrowYes(East).ClearArrowYes(),
			want1: false,
		},
		{
			name:  "WithArrowYes ClearArrowNo",
			c:     Chip(0).WithArrowYes(South).ClearArrowNo(),
			want:  South,
			want1: true,
		},
		{
			name:  "WithArrowYes twice",
			c:     Chip(0).WithArrowYes(South).WithArrowYes(West),
			want:  West,
			want1: true,
		},
		{
			name:  "WithArrowNo",
			c:     Chip(0).WithArrowNo(South),
			want1: false,
		},
		{
			name:  "WithArrowYes WithType",
			c:     Chip(0).WithArrowYes(North).WithType(ForwardChip),
			want:  North,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.ArrowYes()
			if got != tt.want {
				t.Errorf("Chip.ArrowYes() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Chip.ArrowYes() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestChip_ArrowNo(t *testing.T) {
	tests := []struct {
		name  string
		c     Chip
		want  Orientation
		want1 bool
	}{
		{
			name:  "Default",
			c:     Chip(0),
			want1: false,
		},
		{
			name:  "WithArrowNo",
			c:     Chip(0).WithArrowNo(East),
			want:  East,
			want1: true,
		},
		{
			name:  "WithArrowNo ClearArrowNo",
			c:     Chip(0).WithArrowNo(East).ClearArrowNo(),
			want1: false,
		},
		{
			name:  "WithArrowNo ClearArrowYes",
			c:     Chip(0).WithArrowNo(South).ClearArrowYes(),
			want:  South,
			want1: true,
		},
		{
			name:  "WithArrowNo twice",
			c:     Chip(0).WithArrowNo(South).WithArrowNo(West),
			want:  West,
			want1: true,
		},
		{
			name:  "WithArrowYes",
			c:     Chip(0).WithArrowYes(South),
			want1: false,
		},
		{
			name:  "WithArrowNo WithType",
			c:     Chip(0).WithArrowNo(North).WithType(ForwardChip),
			want:  North,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.ArrowNo()
			if got != tt.want {
				t.Errorf("Chip.ArrowNo() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Chip.ArrowNo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
